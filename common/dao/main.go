package dao

import (
	"fmt"
)

const (
	QueryUpdateMetric = "UPDATE"
	QueryInsertMetric = "INSERT"
	QuerySelectMetric = "SELECT"
	QueryDeleteMetric = "DELETE"
	QueryBuildMetric  = "BUILD"
	QueryGetMetric    = "GET"
	QuerySetMetric    = "SET"
	QuerySetExMetric  = "SETEX"
)

var ErrDuplicate = fmt.Errorf("DB duplicate error")
var ErrNoRows = fmt.Errorf("DB no rows in resultset")

type (
	DAOTx interface {
		CommitTx() error
		RollbackTx() error
	}

	ServiceDAO interface {
		BeginTx() (DAOTx, error)
	}
)
