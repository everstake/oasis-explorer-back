package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
	"github.com/wedancedalot/squirrel"
	"oasisTracker/common/baseconf/types"
	"oasisTracker/common/dao"
	"oasisTracker/common/log"
)

const longRunningQueryLimit = time.Millisecond * 100

type Mysql struct {
	Db        *sqlx.DB
	DebugMode bool
}

var errorsRegexp = regexp.MustCompile(`^Error (?P<code>\d+)`)

// CreateConnection creates a mysql connection
func CreateConnection(c *types.DBParams, debugMode bool) (*Mysql, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", c.User, c.Password, c.Host, c.Port, c.Database))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetimeMS) * time.Millisecond)

	return &Mysql{db, debugMode}, nil
}

// Migrate makes migration for DB from migrationsDir
func Migrate(c *types.DBParams, migrationsDir string) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Join(filepath.Dir(ex), migrationsDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		//return errors.New("Migrations dir does not exist: " + dir)
		dir = migrationsDir
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return errors.New("Migrations dir does not exist: " + dir)
		}
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&multiStatements=true&parseTime=true", c.User, c.Password, c.Host, c.Port, c.Database))
	if err != nil {
		return err
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}

	_, err = migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	return err
}

// Find into dest from querybuilder
func (this *Mysql) Find(dest interface{}, b squirrel.SelectBuilder, tx ...*sqlx.Tx) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}

	start := time.Now()

	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].Select(dest, q, params...)
	} else {
		err = this.Db.Select(dest, q, params...)
	}

	exTime := time.Since(start)
	if exTime > longRunningQueryLimit {
		log.Warn("long running query found!:", zap.String("query", q), zap.Any("params", params), zap.Duration("execution time", exTime))
	} else if this.DebugMode {
		log.Debug(q, zap.String("query", q), zap.Duration("execution time", exTime))
	}

	return err
}

// Find into dest from querybuilder
func (this *Mysql) FindRaw(dest interface{}, q string, params ...interface{}) error {
	start := time.Now()
	err := this.Db.Select(dest, q, params...)

	if this.DebugMode {
		log.Debug(q, zap.Any("params", params), zap.Duration("execution time", time.Since(start)))
	}

	return err
}

// Find first row into dest from querybuilder
func (this *Mysql) FindFirst(dest interface{}, b squirrel.SelectBuilder, tx ...*sqlx.Tx) (err error) {
	q, params, err := b.ToSql()
	if err != nil {
		return
	}

	start := time.Now()

	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].Get(dest, q, params...)
	} else {
		err = this.Db.Get(dest, q, params...)
	}

	exTime := time.Since(start)
	if exTime > longRunningQueryLimit {
		log.Warn("long running query found!:", zap.String("query", q), zap.Any("params", params), zap.Duration("execution time", exTime))
	} else if this.DebugMode {
		log.Debug(q, zap.String("query", q), zap.Duration("execution time", exTime))
	}

	return this.parseError(err)
}

// Insert from querybuilder
func (this *Mysql) Insert(b squirrel.InsertBuilder, tx ...*sqlx.Tx) (uint64, error) {
	q, args, err := b.ToSql()
	if err != nil {
		return 0, err
	}

	result, err := this.exec(q, args, tx...)
	if err != nil {
		return 0, this.parseError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), err
}

// Exec query
func (this *Mysql) Exec(q string, args []interface{}, err error, tx ...*sqlx.Tx) (uint64, error) {
	if err != nil {
		return 0, this.parseError(err)
	}

	result, err := this.exec(q, args, tx...)
	if err != nil {
		return 0, this.parseError(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, this.parseError(err)
	}

	return uint64(affectedRows), nil
}

func (this *Mysql) CallFunc(q string, params []interface{}, tx ...*sqlx.Tx) error {
	var err error

	start := time.Now()
	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].QueryRow(q).Scan(params...)
	} else {
		err = this.Db.QueryRow(q).Scan(params...)
	}

	if this.DebugMode {
		log.Debug(q, zap.Duration("execution time", time.Since(start)))
	}

	return err
}

func (this *Mysql) exec(q string, params []interface{}, tx ...*sqlx.Tx) (sql.Result, error) {
	var result sql.Result
	var err error

	start := time.Now()

	if len(tx) > 0 && tx[0] != nil {
		result, err = tx[0].Exec(q, params...)
	} else {
		result, err = this.Db.Exec(q, params...)
	}

	if this.DebugMode {
		log.Debug(q, zap.Any("params", params), zap.Duration("execution time", time.Since(start)))
	}

	return result, err
}

func (this *Mysql) parseError(err error) error {
	if err == nil {
		return nil
	}

	// Just a wrapper not to use sql lib directly from code
	if err == sql.ErrNoRows {
		return dao.ErrNoRows
	}

	matches := this.matchStringGroups(errorsRegexp, err.Error())
	code, ok := matches["code"]
	if !ok {
		return err
	}

	switch code {
	case "1062":
		return dao.ErrDuplicate
	default:
		return err
	}
}

// matchStringGroups matches regexp with capture groups. Returns map string string
func (this *Mysql) matchStringGroups(re *regexp.Regexp, s string) map[string]string {
	m := re.FindStringSubmatch(s)
	n := re.SubexpNames()

	r := make(map[string]string, len(m))

	if len(m) > 0 {
		m, n = m[1:], n[1:]
		for i, _ := range n {
			r[n[i]] = m[i]
		}
	}

	return r
}
