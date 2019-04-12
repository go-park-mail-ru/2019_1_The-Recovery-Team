package postgresql

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"os"
)

// MakeMigrations process database migrations from file
func MakeMigrations(conn *pgx.Conn, filePath string) error {
	tx, err := conn.Begin()
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
