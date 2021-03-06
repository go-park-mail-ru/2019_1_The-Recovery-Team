package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/service"

	"github.com/go-redis/redis"

	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	sessionRepo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 50052, "service port")
	dev := flag.Bool("local", false, "local config flag")
	flag.Parse()

	if *dev {
		viper.SetConfigName("local")
	} else {
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

	redisConn := redis.NewClient(&redis.Options{
		Addr:     redisAddr + ":" + strconv.Itoa(redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := redisConn.Ping().Result(); err != nil {
		log.Fatal("Redis connection refused")
	}
	defer redisConn.Close()

	interactor := usecase.NewSessionInteractor(sessionRepo.NewSessionRepo(redisConn))
	sessionService := session.NewService(interactor)
	server := grpc.NewServer()

	session.RegisterSessionServer(server, sessionService)

	service.RegisterInConsul(consulAddr, consulPort, sessionName, sessionAddr, *port)
	defer service.DeregisterInConsul(consulAddr, consulPort, sessionName, *port)

	log.Print(server.Serve(lis))
}
