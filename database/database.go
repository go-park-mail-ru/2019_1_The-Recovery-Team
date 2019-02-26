package database

import (
	"errors"

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
func (manager *Manager) DB() (*sqlx.DB, error) {
	if manager.db == nil {
		return nil, errors.New("database manager isn't initialized")
	}

	return manager.db, nil
}

// Close closes connection to database
func (manager *Manager) Close() error {
	if manager.db == nil {
		return errors.New("database manager isn't initialized")
	}

	err := manager.db.Close()
	manager.db = nil
	return err
}

// Find returns first data witch satisfies query with args
func (manager *Manager) Find(result interface{}, query string, args ...interface{}) error {
	dbo, err := manager.DB()
	if err != nil {
		return err
	}

	err = dbo.Get(result, query, args...)
	return err
}

// FindAll returns all data witch satisfies query with args
func (manager *Manager) FindAll(result interface{}, query string, args ...interface{}) error {
	dbo, err := manager.DB()
	if err != nil {
		return err
	}

	err = dbo.Select(result, query, args...)
	return err
}
