package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/spf13/viper"

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
	port := flag.Int("port", 50052, "service port")
	dev := flag.Bool("local", false, "local config flag")
	flag.Parse()

	if *dev {
		fmt.Println(*dev)
		viper.SetConfigName("local")
	} else {
		fmt.Println(*dev)
		viper.SetConfigName("config")
	}
	viper.SetConfigType("json")
	viper.AddConfigPath("build/config/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config files:", err)
	}

	consulAddr := viper.GetString("consul.address")
	consulPort := viper.GetInt("consul.port")
	sessionName := viper.GetString("session.name")
	sessionAddr := viper.GetString("session.address")
	redisAddr := viper.GetString("redis.address")
	redisPort := viper.GetInt("redis.port")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("Failed to listen port", port)
	}

	redisConn, err := redis.DialURL(fmt.Sprintf("redis://:@%s:%d", redisAddr, redisPort))
	if err != nil {
		log.Fatal("Redis connection refused")
	}
	defer redisConn.Close()

	interactor := usecase.NewSessionInteractor(sessionRepo.NewSessionRepo(&redisConn))
	service := session.NewService(interactor)
	server := grpc.NewServer()

	session.RegisterSessionServer(server, service)

	config := consulapi.DefaultConfig()
	config.Address = consulAddr + ":" + strconv.Itoa(consulPort)
	consul, err := consulapi.NewClient(config)

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceId + strconv.Itoa(*port),
		Name:    sessionName,
		Port:    *port,
		Address: sessionAddr,
	})
	if err != nil {
		log.Println("Can't add session service to resolver:", err)
		return
	}
	log.Println("Registered in consul", serviceId, port)

	defer func() {
		err := consul.Agent().ServiceDeregister(serviceId + strconv.Itoa(*port))
		if err != nil {
			log.Println("Can't remove service from resolver:", err)
		}
		log.Println("Deregistered in resolver", serviceId, port)
	}()

	log.Print(server.Serve(lis))
}
