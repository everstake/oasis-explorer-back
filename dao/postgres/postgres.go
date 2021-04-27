package postgres

import (
	"database/sql"
	"fmt"
	"oasisTracker/common/dao"
	"oasisTracker/conf"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jinzhu/gorm"
)

const migrationsDir = "./dao/postgres/migrations"

type DAO struct {
	db *gorm.DB
}

func New(c conf.Config) (*DAO, error) {
	d, err := newConn(c)
	if err != nil {
		return nil, err
	}

	err = makeMigration(d.DB(), migrationsDir, c.Postgres.Database, c.Postgres.Schema)
	if err != nil {
		return nil, err
	}

	return &DAO{d}, nil
}

func newConn(c conf.Config) (*gorm.DB, error) {

	db, err := gorm.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Database, c.Postgres.Schema))
	if err != nil {
		return nil, err
	}

	//db.SetLogger(&config.DbLogger{})
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return fmt.Sprintf("%s.%s", c.Postgres.Schema, defaultTableName)
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func makeMigration(conn *sql.DB, migrationDir, dbName, schemaName string) (err error) {
	driver, err := postgres.WithInstance(conn, &postgres.Config{DatabaseName: dbName, SchemaName: schemaName})
	if err != nil {
		return err
	}

	err = dao.MakeMigration(driver, migrationDir, dbName)
	if err != nil {
		return err
	}

	return nil
}
