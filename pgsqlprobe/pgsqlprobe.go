package pgsqlprobe

import (
	"context"
	"database/sql"

	// Mandatory for sql.Open to work
	_ "github.com/lib/pq"
	errgo "gopkg.in/errgo.v1"
)

type PostgreSQLProbe struct {
	name             string
	connectionString string
}

// NewPostgreSQLProbe instantiate a new PostgreSQL probe:
// - name: probe name
// - connectionString: connection string with the form "postgres://username:password@example.com"
func NewPostgreSQLProbe(name, connectionString string) PostgreSQLProbe {
	return PostgreSQLProbe{
		name:             name,
		connectionString: connectionString,
	}
}

func (p PostgreSQLProbe) Name() string {
	return p.name
}

func (p PostgreSQLProbe) Check(_ context.Context) error {
	client, err := sql.Open("postgres", p.connectionString)
	if err != nil {
		return errgo.Notef(err, "fail to open a new connection to PostgreSQL")
	}
	defer client.Close()

	err = client.Ping()
	if err != nil {
		return errgo.Notef(err, "unable to contact PostgreSQL host")
	}

	return nil
}
