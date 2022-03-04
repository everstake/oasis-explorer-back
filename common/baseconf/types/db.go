package types

import "fmt"

// MysqlParams is mysql config struct
type DBParams struct {
	Host              string
	Port              int
	User              string
	Password          string
	Database          string
	Schema            string
	DebugMode         bool
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetimeMS uint64
}

// Validate checks all MysqlParams fields
func (this *DBParams) Validate() error {
	if this.Port == 0 {
		return fmt.Errorf("bad db Port")
	}

	if this.Host == "" {
		return fmt.Errorf("no db Host")
	}

	if this.User == "" {
		return fmt.Errorf("no db User")
	}

	// if this.Password == "" {
	// 	return fmt.Errorf("no mysql Password")
	// }

	if this.Database == "" {
		return fmt.Errorf("no Database")
	}

	return nil
}
