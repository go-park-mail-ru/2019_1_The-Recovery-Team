package main

import (
	"log"
	"net"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql"
	repo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

const (
	serviceId = "SProfile_127.0.0.1:"
)

func main() {
	port := os.Getenv("PORT")
	portI, err := strconv.Atoi(port)
	if err != nil {
		port = "50051"
		portI = 50051
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen port", port)
	}

	psqlConfig := pgx.ConnConfig{
		Host:     "db",
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

	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := consulapi.NewClient(config)

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceId + port,
		Name:    "profile-service",
		Port:    portI,
		Address: "profile",
	})
	if err != nil {
		log.Println("Can't add profile service to resolver:", err)
		return
	}
	log.Println("Registered in resolver", serviceId, port)

	defer func() {
		err := consul.Agent().ServiceDeregister(serviceId + port)
		if err != nil {
			log.Println("Can't remove service from resolver:", err)
		}
		log.Println("Deregistered in resolver", serviceId, port)
	}()

	log.Print(server.Serve(lis))
}
