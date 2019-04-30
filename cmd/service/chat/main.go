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
	chatDbAddr := viper.GetString("chat.database.address")
	chatDbPort := viper.GetInt("chat.database.port")
	chatDbMigrationsFile := viper.GetString("chat.database.migrations.file")
	consulAddr := viper.Get("consul.address")
	sessionName := viper.Get("session.name")

	chatConfig := pgx.ConnConfig{
		Host:     chatDbAddr,
		Port:     uint16(chatDbPort),
		Database: "sadislandschat",
		User:     "recoveryteam",
		Password: "123456",
	}

	chatConn, err := pgx.Connect(chatConfig)
	if err != nil {
		log.Fatal("Chat postresql connection refused")
	}

	if err := postgresql.MakeMigrations(chatConn, chatDbMigrationsFile); err != nil {
		log.Fatal("Database migrations failed", err)
	}
	pgxClose(chatConn)

	chatConn, err = pgx.Connect(chatConfig)
	if err != nil {
		log.Fatal("Chat postresql connection refused")
	}
	defer pgxClose(chatConn)

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

	messageInteractor := usecase.NewMessageInteractor(message.NewRepo(chatConn))
	chatInteractor := usecase.NewChatInteractor(chat.NewRepo(logger, messageInteractor))
	sessionManager := session.NewSessionClient(sessionConn)

	api := chatApi.NewApi(chatInteractor, &sessionManager, logger)
	log.Print(http.ListenAndServe(":"+strconv.Itoa(port), api.Router))
}
