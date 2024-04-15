package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	// PostgreSQL
	_ "github.com/lib/pq"
)

type (
	// db is wrapper for Master and Slave database connection
	Replication struct {
		DriverName string
		Master     *sqlx.DB
		Slave      *sqlx.DB
	}

	// ConnectionOptions list of option to connect to database
	ConnectionOptions struct {
		Retry                 int
		MaxOpenConnections    int
		MaxIdleConnections    int
		ConnectionMaxLifetime time.Duration
	}
)

func (db *Replication) Primary() *sqlx.DB {
	return db.Master
}

func (db *Replication) Secondary() *sqlx.DB {
	if db.Slave != nil {
		return db.Slave
	}
	return db.Master
}

// Get return one value in destination interface.
// It will return error when no value returned.
func (db *Replication) Get(destination interface{}, query string, args ...interface{}) error {
	return db.Secondary().Get(destination, query, args...)
}

// Select return more than one value in destination using reflection.
func (db *Replication) Select(destination interface{}, query string, args ...interface{}) error {
	return db.Secondary().Select(destination, query, args...)
}

// SelectLeader return more than one value in destination using reflection.
func (db *Replication) SelectLeader(destination interface{}, query string, args ...interface{}) error {
	return db.Master.Select(destination, query, args...)
}

// Query database and return *sql.Rows
func (db *Replication) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Secondary().Query(query, args...)
}

// Queryx queries the database and returns an *sqlx.Rows.
// Any placeholder parameters are replaced with supplied args.
func (db *Replication) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.Secondary().Queryx(query, args...)
}

// QueryRow expecting to return at least one *sql.Row
func (db *Replication) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.Secondary().QueryRow(query, args...)
}

// QueryRowx expecting to return at least one row
func (db *Replication) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return db.Secondary().QueryRowx(query, args...)
}

// NamedQuery using this db.
// Any named placeholder parameters are replaced with fields from arg.
func (db *Replication) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return db.Secondary().NamedQuery(query, arg)
}

// Exec executes query without returning rows.
func (db *Replication) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Master.Exec(query, args...)
}

// NamedExec uses BindStruct to get a query executable by the driver and
// then runs Exec on the result.  Returns an error from the binding
// or the query excution itself.
func (db *Replication) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.Master.NamedExec(query, arg)
}

// Begin return sql transaction object
func (db *Replication) Begin() (*sql.Tx, error) {
	return db.Master.Begin()
}

// Beginx return sqlx transaction object
func (db *Replication) Beginx() (*sqlx.Tx, error) {
	return db.Master.Beginx()
}

// Rebind a query to targeted bind type
func (db *Replication) Rebind(query string) string {
	return sqlx.Rebind(sqlx.BindType(db.DriverName), query)
}

// Named takes a query using named parameters and an argument and
// returns a new query with a list of args that can be executed by a database.
func (db *Replication) Named(query string, arg interface{}) (string, interface{}, error) {
	return sqlx.Named(query, arg)
}

// PrepareNamedContextLeader returns an sqlx.NamedStmt
func (db *Replication) PrepareNamedContextLeader(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return db.Master.PrepareNamedContext(ctx, db.Master.Rebind(query))
}

// PrepareNamedContextFollower returns an sqlx.NamedStmt
func (db *Replication) PrepareNamedContextFollower(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	return db.Secondary().PrepareNamedContext(ctx, db.Slave.Rebind(query))
}

// PreparexContextFollower returns an sqlx.Stmt
func (db *Replication) PreparexContextFollower(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return db.Secondary().PreparexContext(ctx, query)
}
