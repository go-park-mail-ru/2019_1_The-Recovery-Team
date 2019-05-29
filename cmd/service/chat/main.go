package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/service"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/database"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	chatApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/chat"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/message"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func pgxClose(conn *pgx.Conn) {
	err := conn.Close()
	if err != nil {
		log.Fatal("pgx connection close failed", err)
	}
}

func main() {
	dev := flag.Bool("local", false, "local config flag")
	dbUser := flag.String("db_user", "recoveryteam", "database username")
	dbPassword := flag.String("db_password", "123456", "database password")
	dbName := flag.String("db_name", "sadislandschat", "database name")
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
	port := viper.GetInt("chat.port")
	postgresqlAddr := viper.GetString("chat.database.address")
	postgresqlPort := viper.GetInt("chat.database.port")
	migrationsFile := viper.GetString("chat.database.migrations.file")
	sessionName := viper.GetString("session.name")

	// Start watching for updates in consul
	resolver.RegisterDefault(consulAddr, consulPort, 5*time.Second)

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

	sessionConn := service.Connect(consulAddr, sessionName)
	defer sessionConn.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	messageInteractor := usecase.NewMessageInteractor(message.NewRepo(dbConnPool))
	chatInteractor := usecase.NewChatInteractor(chat.NewRepo(logger, messageInteractor))
	sessionManager := session.NewSessionClient(sessionConn)

	go chatInteractor.Run()

	api := chatApi.NewApi(chatInteractor, &sessionManager, logger)
	logger.Info("Chat exit",
		zap.Error(http.ListenAndServe(":"+strconv.Itoa(port), api.Router)))
}
