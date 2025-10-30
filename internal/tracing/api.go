package tracing

import (
	"context"

	"github.com/jmoiron/sqlx"

	"go.elastic.co/apm/module/apmsql/v2"
	"go.elastic.co/apm/v2"
)

func Open(driver string, dsn string) (*sqlx.DB, error) {
	conn, err := apmsql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return sqlx.NewDb(conn, driver), nil
}

func Connect(driver string, dsn string) (*sqlx.DB, error) {
	conn, err := Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, nil
}

func CaptureError(ctx context.Context, err error) {
	apm.CaptureError(ctx, err).Send()
}
