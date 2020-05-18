package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	goclickhouse "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mailru/go-clickhouse"
	"oasisTracker/conf"
	"oasisTracker/dao/clickhouse/client"
	"strings"
)

type Clickhouse struct {
	db *DB
}

const migrationsDir = "./dao/clickhouse/migrations"

func New(cfg conf.Config) (db *Clickhouse, err error) {
	config := cfg.Clickhouse
	config.Database = cfg.Scanner.Database
	conn, err := newConnection(config)
	if err != nil {
		return db, fmt.Errorf("newConnection: %s", err.Error())
	}

	client := NewDB(conn)

	return &Clickhouse{
		db: client,
	}, nil
}

func (c *Clickhouse) GetChain() interface{} {
	return c.db
}

func newConnection(cfg conf.Clickhouse) (*sql.DB, error) {
	conn, err := sql.Open("clickhouse", makeSource(cfg))
	if err != nil {
		return nil, fmt.Errorf("can`t make connection: %s", err.Error())
	}

	err = makeMigration(conn, migrationsDir, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("can`t make makeMigration: %s", err.Error())
	}
	return conn, nil
}

func makeSource(cfg conf.Clickhouse) string {
	return fmt.Sprintf("%s://%s:%d/%s?password=%s&user=%s",
		strings.Trim(cfg.Protocol, "://"),
		strings.Trim(cfg.Host, "/"),
		cfg.Port,
		cfg.Database,
		cfg.Password,
		cfg.User,
	)
}

func makeMigration(conn *sql.DB, migrationDir string, dbName string) error {
	driver, err := goclickhouse.WithInstance(conn, &goclickhouse.Config{})
	if err != nil {
		return fmt.Errorf("clickhouse.WithInstance: %s", err.Error())
	}
	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		dbName, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %s", err.Error())
	}
	if err := mg.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}

type (
	DB struct {
		conn   *sql.DB
		client client.Client
	}
)

func NewDB(conn *sql.DB) *DB {
	return &DB{
		conn:   conn,
		client: client.New(conn),
	}
}
