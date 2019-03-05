package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

// Manager is a wrapper around sqlx.DB
type Manager struct {
	db *sqlx.DB
}

// InitDatabaseManager initialize manager with connection to database
func InitDatabaseManager(username, password, host, name string) (*Manager, error) {
	var err error
	manager := &Manager{}
	manager.db, err = sqlx.Open("postgres", "postgres://"+username+":"+password+"@"+host+"/"+name+"?sslmode=disable")
	if err != nil {
		return nil, err
	}

	if err := manager.db.Ping(); err != nil {
		return nil, err
	}

	return manager, nil
}

// DB returns connection to database
func (manager *Manager) DB() *sqlx.DB {
	return manager.db
}

// Close closes connection to database
func (manager *Manager) Close() error {
	if manager.db == nil {
		return nil
	}
	err := manager.db.Close()
	manager.db = nil
	return err
}

// Find returns first data witch satisfies query with args
func (manager *Manager) Find(result interface{}, query string, args ...interface{}) error {
	dbo := manager.DB()
	err := dbo.Get(result, query, args...)
	return err
}

// FindAll returns all data witch satisfies query with args
func (manager *Manager) FindAll(result interface{}, query string, args ...interface{}) error {
	dbo := manager.DB()
	err := dbo.Select(result, query, args...)
	return err
}

// FindWithField checks existence of row with field value
func (manager *Manager) FindWithField(table, field, value string) (bool, error) {
	var result string
	dbo := manager.DB()
	query := `SELECT ` + field + ` FROM ` + table + ` WHERE LOWER(` + field + `)` + ` = LOWER($1)`
	err := dbo.Get(&result, query, value)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Create adds a new entry
func (manager *Manager) Create(result interface{}, query string, args ...interface{}) error {
	dbo := manager.DB()
	row := dbo.QueryRowx(query, args...)
	if row.Err() != nil {
		return row.Err()
	}

	err := row.StructScan(result)
	if err != nil {
		return err
	}
	return nil
}
