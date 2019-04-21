package main

import (
	"log"
	"net"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql"
	repo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("Failed to listen port 9091")
	}

	psqlConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "sadislands",
		User:     "recoveryteam",
		Password: "123456",
	}

	// Create connection for migrations
	psqlConn, err := pgx.Connect(psqlConfig)
	if err != nil {
		log.Fatal("Postgresql connection refused")
	}

	if err := postgresql.MakeMigrations(psqlConn, "build/schema/0_initial.sql"); err != nil {
		log.Fatal("Database migrations failed", err)
	}
	psqlConn.Close()

	// Create new connection to database with updated OIDs
	psqlConn, err = pgx.Connect(psqlConfig)
	if err != nil {
		log.Fatal("Postgresql connection refused")
	}
	defer psqlConn.Close()

	interactor := usecase.NewProfileInteractor(repo.NewRepo(psqlConn))
	service := profile.NewService(interactor)
	server := grpc.NewServer()

	profile.RegisterProfileServer(server, service)
	log.Print(server.Serve(lis))
}
