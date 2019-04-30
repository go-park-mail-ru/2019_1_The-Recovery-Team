package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/message"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	chatApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func pgxClose(conn *pgx.Conn) {
	err := conn.Close()
	if err != nil {
		log.Println("pgx connection close failed", err)
	}
}

func init() {
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

	addr := viper.GetString("consul.address")
	port := viper.GetInt("consul.port")
	resolver.RegisterDefault(addr, port, 5*time.Second)
}

func main() {
	port := viper.GetInt("chat.port")
	postgresqlAddr := viper.GetString("chat.database.address")
	postgresqlPort := viper.GetInt("chat.database.port")
	migrationsFile := viper.GetString("chat.database.migrations.file")
	consulAddr := viper.Get("consul.address")
	sessionName := viper.Get("session.name")

	psqlConfig := pgx.ConnConfig{
		Host:     postgresqlAddr,
		Port:     uint16(postgresqlPort),
		Database: "sadislandschat",
		User:     "recoveryteam",
		Password: "123456",
	}

	psqlConfigPool := pgx.ConnPoolConfig{
		ConnConfig:     psqlConfig,
		MaxConnections: 50,
	}

	// Create connection for migrations
	psqlConn, err := pgx.Connect(psqlConfig)
	if err != nil {
		log.Fatal("Postresql connection refused")
	}

	if err := postgresql.MakeMigrations(psqlConn, migrationsFile); err != nil {
		log.Fatal("Database migrations failed", err)
	}
	pgxClose(psqlConn)

	// Create new connection to database with updated OIDs
	psqlConnPool, err := pgx.NewConnPool(psqlConfigPool)
	if err != nil {
		log.Fatal("Postresql connection refused")
	}
	defer psqlConnPool.Close()

	sessionConn, err := grpc.Dial(fmt.Sprintf("srv://%s/%s", consulAddr, sessionName),
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Authentication service connection refused:", err)
	}
	defer sessionConn.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	messageInteractor := usecase.NewMessageInteractor(message.NewRepo(psqlConnPool))
	chatInteractor := usecase.NewChatInteractor(chat.NewRepo(logger, messageInteractor))
	sessionManager := session.NewSessionClient(sessionConn)

	api := chatApi.NewApi(chatInteractor, &sessionManager, logger)
	log.Print(http.ListenAndServe(":"+strconv.Itoa(port), api.Router))
}
