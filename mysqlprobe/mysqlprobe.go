package mysqlprobe

import (
	"context"
	"database/sql"

	// Mandatory for sql.Open to work
	_ "github.com/go-sql-driver/mysql"

	errgo "gopkg.in/errgo.v1"
)

type MySQLProbe struct {
	name             string
	connectionString string
}

// NewMySQLProbe instantiate a new MySQL probe:
// - name: probe name
// - connectionString: connection string with the form "mysql://username:password@example.com"
func NewMySQLProbe(name, connectionString string) MySQLProbe {
	return MySQLProbe{
		name:             name,
		connectionString: connectionString,
	}
}

func (p MySQLProbe) Name() string {
	return p.name
}

func (p MySQLProbe) Check(_ context.Context) error {
	client, err := sql.Open("mysql", p.connectionString)
	if err != nil {
		return errgo.Notef(err, "fail to open a new connection to MySQL")
	}
	defer client.Close()

	err = client.Ping()
	if err != nil {
		return errgo.Notef(err, "unable to contact MySQL host")
	}

	return nil
}
