package main

import (
	"log"
	"net"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/session"

	repo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/usecase"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to listen port 8081")
	}

	redisConn, err := redis.DialURL("redis://:@localhost:6379")
	if err != nil {
		log.Fatal("Redis connection refused")
	}
	defer redisConn.Close()

	interactor := usecase.NewSessionInteractor(repo.NewSessionRepo(&redisConn))
	service := session.NewService(interactor)
	server := grpc.NewServer()

	session.RegisterSessionServer(server, service)
	log.Print(server.Serve(lis))
}
