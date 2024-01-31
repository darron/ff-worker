package d1

import (
	"database/sql"
	"database/sql/driver"
	"errors"
)

func init() {
	sql.Register("d1", &Driver{})
}

type Driver struct{}

var (
	_ driver.Driver = (*Driver)(nil)
)

func (d *Driver) Open(string) (driver.Conn, error) {
	return nil, errors.New("d1: Open is not supported. use d1.OpenConnector and sql.OpenDB instead")
}
