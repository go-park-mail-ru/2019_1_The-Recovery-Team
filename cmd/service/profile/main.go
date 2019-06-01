package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/service"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/database"

	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	profileRepo "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"github.com/jackc/pgx"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 50051, "service port")
	dev := flag.Bool("local", false, "local config flag")
	dbUser := flag.String("db_user", "recoveryteam", "database username")
	dbPassword := flag.String("db_password", "123456", "database password")
	dbName := flag.String("db_name", "sadislands", "database name")
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

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("Failed to listen port", port)
	}

	dbConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     postgresqlAddr,
			Port:     uint16(postgresqlPort),
			Database: *dbName,
			User:     *dbUser,
			Password: *dbPassword,
		},
		MaxConnections: 50,
	}
	dbConnPool := database.Connect(dbConfig, migrationsFile)
	defer dbConnPool.Close()

	interactor := usecase.NewProfileInteractor(profileRepo.NewRepo(dbConnPool))
	profileService := profile.NewService(interactor)
	server := grpc.NewServer()

	profile.RegisterProfileServer(server, profileService)

	service.RegisterInConsul(consulAddr, consulPort, profileName, profileAddr, *port)
	defer service.DeregisterInConsul(consulAddr, consulPort, profileName, *port)

	log.Print(server.Serve(lis))
}
