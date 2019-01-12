package sql

import (
	"context"
	stdlibsql "database/sql"
	"time"

	"code.cloudfoundry.org/lager"
)

// DB is a dumb wrapper around a stdlib sql.DB that also logs queries.
// All of its functions simply log what they are doing and then call down
// to the corresponding stdlib sql.DB function.
type DB struct {
	logger lager.Logger

	db *stdlibsql.DB
}

// Open creates a DB via a driverName and a dataSourceName. It calls the stdlib
// sql.Open function.
func Open(logger lager.Logger, driverName, dataSourceName string) (*DB, error) {
	db, err := stdlibsql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{logger: logger, db: db}, nil
}

func (db *DB) Exec(query string, args ...interface{}) (stdlibsql.Result, error) {
	db.logger.Debug("exec", lager.Data{"query": query, "args": args})
	if args == nil {
		return db.db.ExecContext(makeCtx(), query)
	} else {
		return db.db.ExecContext(makeCtx(), query, args)
	}
}

func (db *DB) Query(query string, args ...interface{}) (*stdlibsql.Rows, error) {
	db.logger.Debug("query", lager.Data{"query": query, "args": args})
	if args == nil {
		return db.db.QueryContext(makeCtx(), query)
	} else {
		return db.db.QueryContext(makeCtx(), query, args)
	}
}

func (db *DB) QueryRow(query string, args ...interface{}) *stdlibsql.Row {
	db.logger.Debug("query", lager.Data{"query-row": query, "args": args})
	if args == nil {
		return db.db.QueryRowContext(makeCtx(), query)
	} else {
		return db.db.QueryRowContext(makeCtx(), query, args...)
	}
}

func (db *DB) Close() error {
	db.logger.Debug("close")
	return db.db.Close()
}

func (db *DB) Prepare(query string) (*stmt, error) {
	db.logger.Debug("prepare", lager.Data{"query": query})

	s, err := db.db.PrepareContext(makeCtx(), query)
	if err != nil {
		return nil, err
	}

	return &stmt{logger: db.logger.Session("prepare"), stmt: s}, nil
}

type stmt struct {
	logger lager.Logger

	stmt *stdlibsql.Stmt
}

func (s *stmt) Exec(args ...interface{}) (stdlibsql.Result, error) {
	s.logger.Debug("exec", lager.Data{"args": args})
	return s.stmt.ExecContext(makeCtx(), args...)
}

func (s *stmt) Close() error {
	s.logger.Debug("close")
	return s.stmt.Close()
}

func makeCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return ctx
}
