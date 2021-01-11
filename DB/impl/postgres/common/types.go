package common

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"time"
)

const (
	DEFAULT_REQUEST_TTL = 2 * time.Second
)

/*
	these types are borrowed from the PGX source code and exposed here to facilitate easier method chaining
*/

type (
	Execer interface {
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	}

	Queryer interface {
		Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	}

	QueryRower interface {
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	}
)
