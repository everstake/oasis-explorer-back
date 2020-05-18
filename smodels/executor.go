package smodels

type Executor struct {
	ExecHeight func(uint64) error
	Save       func() error
}
