package smodels

import "google.golang.org/grpc"

type Executor struct {
	ExecHeight func(*grpc.ClientConn, uint64) error
	Truncate   func()
	Save       func() error
}
