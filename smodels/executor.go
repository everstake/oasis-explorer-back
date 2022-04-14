package smodels

type Executor struct {
	ExecHeight func(uint64) error
	Truncate   func()
	Save       func() error
}
