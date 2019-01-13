package sql

import (
	"context"
	stdlibsql "database/sql"
	"time"

	"code.cloudfoundry.org/lager"
)

// DB is a dumb wrapper around a stdlib sql.DB.
//
// All of its functions simply log what they are doing and then call down
// to the corresponding stdlib sql.DB function.
type DB struct {
	db *stdlibsql.DB
}

// Open creates a DB via a driverName and a dataSourceName. It calls the stdlib
// sql.Open function.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := stdlibsql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Second)

	return &DB{db: db}, nil
}

func (db *DB) Exec(
	ctx context.Context,
	logger lager.Logger,
	query string,
	args ...interface{},
) (stdlibsql.Result, error) {
	logger.Debug("exec", lager.Data{"query": query, "args": args})

	if args == nil {
		return db.db.ExecContext(ctx, query)
	} else {
		return db.db.ExecContext(ctx, query, args)
	}
}

func (db *DB) Query(
	ctx context.Context,
	logger lager.Logger,
	query string,
	args ...interface{},
) (*stdlibsql.Rows, error) {
	logger.Debug("query", lager.Data{"query": query, "args": args})

	if args == nil {
		return db.db.QueryContext(ctx, query)
	} else {
		return db.db.QueryContext(ctx, query, args)
	}
}

func (db *DB) QueryRow(
	ctx context.Context,
	logger lager.Logger,
	query string,
	args ...interface{},
) *stdlibsql.Row {
	logger.Debug("query", lager.Data{"query-row": query, "args": args})

	if args == nil {
		return db.db.QueryRowContext(ctx, query)
	} else {
		return db.db.QueryRowContext(ctx, query, args...)
	}
}

func (db *DB) Close(logger lager.Logger) error {
	logger.Debug("close")
	return db.db.Close()
}

func (db *DB) Prepare(
	ctx context.Context,
	logger lager.Logger,
	query string,
) (*stmt, error) {
	logger.Debug("prepare", lager.Data{"query": query})

	s, err := db.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &stmt{stmt: s}, nil
}

type stmt struct {
	stmt *stdlibsql.Stmt
}

func (s *stmt) Exec(
	ctx context.Context,
	logger lager.Logger,
	args ...interface{},
) (stdlibsql.Result, error) {
	logger.Debug("exec", lager.Data{"args": args})

	return s.stmt.ExecContext(ctx, args...)
}

func (s *stmt) Close(logger lager.Logger) error {
	logger.Debug("close")
	return s.stmt.Close()
}
