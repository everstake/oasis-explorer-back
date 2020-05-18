package types

import "fmt"

// MysqlParams is mysql config struct
type MysqlParams struct {
	Host              string
	Port              int
	User              string
	Password          string
	Database          string
	DebugMode         bool
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetimeMS uint64
}

// Validate checks all MysqlParams fields
func (this MysqlParams) Validate() error {
	if this.Port == 0 {
		return fmt.Errorf("bad mysql Port")
	}

	if this.Host == "" {
		return fmt.Errorf("no mysql Host")
	}

	if this.User == "" {
		return fmt.Errorf("no mysql User")
	}

	// if this.Password == "" {
	// 	return fmt.Errorf("no mysql Password")
	// }

	if this.Database == "" {
		return fmt.Errorf("no mysql Database")
	}

	return nil
}
