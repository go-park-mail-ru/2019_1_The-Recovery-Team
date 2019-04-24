package main

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	sessionRepo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

const (
	serviceId = "SSession_127.0.0.1:"
)

func main() {
	port := os.Getenv("PORT")
	portI, err := strconv.Atoi(port)
	if err != nil {
		port = "50052"
		portI = 50052
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen port", port)
	}

	redisConn, err := redis.DialURL("redis://:@redis:6379")
	if err != nil {
		log.Fatal("Redis connection refused")
	}
	defer redisConn.Close()

	interactor := usecase.NewSessionInteractor(sessionRepo.NewSessionRepo(&redisConn))
	service := session.NewService(interactor)
	server := grpc.NewServer()

	session.RegisterSessionServer(server, service)

	config := consulapi.DefaultConfig()
	config.Address = "consul:8500"
	consul, err := consulapi.NewClient(config)

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceId + port,
		Name:    "session-service",
		Port:    portI,
		Address: "session",
	})
	if err != nil {
		log.Println("Can't add session service to resolver:", err)
		return
	}
	log.Println("Registered in consul", serviceId, port)

	defer func() {
		err := consul.Agent().ServiceDeregister(serviceId + port)
		if err != nil {
			log.Println("Can't remove service from resolver:", err)
		}
		log.Println("Deregistered in resolver", serviceId, port)
	}()

	log.Print(server.Serve(lis))
}
