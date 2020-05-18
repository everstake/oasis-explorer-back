package mysql

import (
	"oasisTracker/common/log"
	"oasisTracker/common/mysql"
	"oasisTracker/conf"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"os"
)

const defaultMigrationsDir = "./dao/mysql/migrations"

type (
	DAOTx interface {
		Commit() error
		Rollback()
	}
)

type MysqlDAO struct {
	mysql *mysql.Mysql
}

type mysqlDAOTx struct {
	tx       *sqlx.Tx
	commited bool
}

func NewMysqlConnection(c conf.Config) (*mysql.Mysql, error) {
	m, err := mysql.CreateConnection(&c.Mysql, c.Mysql.DebugMode)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func New(c conf.Config) (*MysqlDAO, error) {
	m, err := NewMysqlConnection(c)
	if err != nil {
		return nil, err
	}
	migrationsDir := defaultMigrationsDir
	if migrations := os.Getenv("MYSQL_MIGRATIONS_PATH"); migrations != "" {
		migrationsDir = migrations
	}
	err = mysql.Migrate(&c.Mysql, migrationsDir)
	if err != nil {
		return nil, err
	}

	return &MysqlDAO{mysql: m}, nil
}

func (md *MysqlDAO) BeginTx() (DAOTx, error) {
	tx, err := md.mysql.Db.Beginx()
	if err != nil {
		return nil, err
	}

	return &mysqlDAOTx{tx: tx}, nil
}

func (md *mysqlDAOTx) Commit() error {
	if md.tx == nil {
		return fmt.Errorf("Tx not initialized")
	}
	err := md.tx.Commit()
	if err != nil {
		return err
	}
	md.commited = true
	return nil
}

func (md *mysqlDAOTx) Rollback() {
	if md.tx == nil {
		return
	}
	err := md.tx.Rollback()
	if err != nil && !md.commited {
		log.Error("error while rolling back tx: ", zap.Error(err))
	}
}

// getTx checks if tx param is not empty and returns an *sqlx.Tx
func daoTx2Sqlx(tx DAOTx) *sqlx.Tx {
	if tx == nil {
		return nil
	}

	t, ok := tx.(*mysqlDAOTx)
	if !ok {
		return nil
	}

	return t.tx
}
