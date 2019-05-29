package database

import (
	"log"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"
	"github.com/jackc/pgx"
)

func Connect(config pgx.ConnPoolConfig, migrationsFile string) *pgx.ConnPool {
	if migrationsFile != "" {
		// Create connection for migrations
		conn, err := pgx.Connect(config.ConnConfig)
		if err != nil {
			log.Fatal("Postresql connection refused")
		}

		if err := postgresql.MakeMigrations(conn, migrationsFile); err != nil {
			log.Fatal("Database migrations failed:", err)
		}
		conn.Close()
	}

	// Create new connection to database with updated OIDs
	connPool, err := pgx.NewConnPool(config)
	if err != nil {
		log.Fatal("Postresql connection refused")
	}
	return connPool
}
