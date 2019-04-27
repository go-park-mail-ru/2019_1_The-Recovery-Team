package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"
	profileRepo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	consulapi "github.com/hashicorp/consul/api"

	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

const (
	serviceId = "SProfile_"
)

func pgxClose(conn *pgx.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println("pgx connection close failed", err)
	}
}

func main() {
	port := flag.Int("port", 50051, "service port")
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
	profileName := viper.GetString("profile.name")
	profileAddr := viper.GetString("profile.address")
	postgresqlPort := viper.GetInt("postgresql.port")
	postgresqlAddr := viper.GetString("postgresql.address")
	migrationsFile := viper.GetString("postgresql.migrations.file")

	chatDbAddr := viper.GetString("chat.database.address")
	chatDbPort := viper.GetInt("chat.database.port")
	chatDbMigrationsFile := viper.GetString("chat.database.migrations.file")

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("Failed to listen port", port)
	}

	psqlConfig := pgx.ConnConfig{
		Host:     postgresqlAddr,
		Port:     uint16(postgresqlPort),
		Database: "sadislands", // sadislandschat
		User:     "recoveryteam",
		Password: "123456",
	}

	chatConfig := pgx.ConnConfig{
		Host:     chatDbAddr,
		Port:     uint16(chatDbPort),
		Database: "sadislandschat",
		User:     "recoveryteam",
		Password: "123456",
	}

	// Create connection for migrations
	psqlConn, err := pgx.Connect(psqlConfig)
	if err != nil {
		log.Fatal("Postgresql connection refused")
	}

	chatConn, err := pgx.Connect(chatConfig)
	if err != nil {
		log.Fatal("Chat postresql connection refused")
	}

	if err := postgresql.MakeMigrations(psqlConn, migrationsFile); err != nil {
		log.Fatal("Database migrations failed", err)
	}
	pgxClose(psqlConn)

	if err := postgresql.MakeMigrations(chatConn, chatDbMigrationsFile); err != nil {
		log.Fatal("Database migrations failed", err)
	}
	pgxClose(chatConn)

	// Create new connection to database with updated OIDs
	psqlConn, err = pgx.Connect(psqlConfig)
	if err != nil {
		log.Fatal("Postgresql connection refused")
	}
	defer pgxClose(psqlConn)

	chatConn, err = pgx.Connect(chatConfig)
	if err != nil {
		log.Fatal("Chat postresql connection refused")
	}
	defer pgxClose(chatConn)

	interactor := usecase.NewProfileInteractor(profileRepo.NewRepo(psqlConn))
	service := profile.NewService(interactor)
	server := grpc.NewServer()

	profile.RegisterProfileServer(server, service)

	config := consulapi.DefaultConfig()
	config.Address = consulAddr + ":" + strconv.Itoa(consulPort)
	consul, err := consulapi.NewClient(config)

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceId + strconv.Itoa(*port),
		Name:    profileName,
		Port:    *port,
		Address: profileAddr,
	})
	if err != nil {
		log.Println("Can't add profile service to resolver:", err)
		return
	}
	log.Println("Registered in resolver", serviceId, port)

	defer func() {
		err := consul.Agent().ServiceDeregister(serviceId + strconv.Itoa(*port))
		if err != nil {
			log.Println("Can't remove service from resolver:", err)
		}
		log.Println("Deregistered in resolver", serviceId, port)
	}()

	log.Print(server.Serve(lis))
}
