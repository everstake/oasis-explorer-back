package dao

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MysqlDAOTx struct {
	tx *sqlx.Tx
}

func (this *MysqlDAOTx) SetTx(dbTx *sqlx.Tx) {
	this.tx = dbTx
}

func (this *MysqlDAOTx) CommitTx() error {
	if this.tx == nil {
		return fmt.Errorf("tx not initialized")
	}

	return this.tx.Commit()
}

func (this *MysqlDAOTx) RollbackTx() error {
	if this.tx == nil {
		return nil
	}

	return this.tx.Rollback()
}

func DaoTx2Sqlx(tx DAOTx) *sqlx.Tx {
	if tx == nil {
		return nil
	}

	t, ok := tx.(*MysqlDAOTx)
	if !ok {
		return nil
	}

	return t.tx
}
