package database

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq" // postgres driver
)

var config = pgx.ConnConfig{
	Host:     "db",
	Port:     5432,
	Database: "sadislands",
	User:     "recoveryteam",
	Password: "123456",
}

// Manager is a wrapper around sqlx.DB
type Manager struct {
	conn *pgx.Conn
}

// InitDatabaseManager initialize manager with connection to database
func InitDatabaseManager(migrationsFilePath string) (*Manager, error) {
	var err error
	manager := &Manager{}
	manager.conn, err = pgx.Connect(config)
	if err != nil {
		return nil, err
	}

	if !manager.conn.IsAlive() {
		return nil, ErrConnRefused
	}

	err = manager.migrate(migrationsFilePath)
	if err != nil {
		return nil, err
	}

	return manager, nil
}

func (dbm *Manager) migrate(filePath string) error {
	if dbm.conn == nil {
		return errors.New("connection doesn't exist")
	}

	tx, err := dbm.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	migrations, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(migrations))
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

// DB returns connection to database
func (dbm *Manager) Conn() *pgx.Conn {
	return dbm.conn
}

// Close closes connection to database
func (dbm *Manager) Close() error {
	if dbm.conn == nil {
		return nil
	}
	err := dbm.conn.Close()
	dbm.conn = nil
	return err
}
