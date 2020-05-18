package client

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/wedancedalot/squirrel"
)

type Client struct {
	conn *sqlx.DB
}

func New(conn *sql.DB) Client {
	return Client{
		conn: sqlx.NewDb(conn, "clickhouse"),
	}
}

func (cli *Client) Find(dest interface{}, b squirrel.SelectBuilder) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}
	err = cli.conn.Select(dest, q, params...)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) FindFirst(dest interface{}, b squirrel.SelectBuilder) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}
	err = cli.conn.Get(dest, q, params...)
	if err != nil {
		return err
	}
	return nil
}

func (cli *Client) Exec(b squirrel.InsertBuilder) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}
	_, err = cli.conn.Exec(q, params...)
	if err != nil {
		return err
	}
	return nil
}
