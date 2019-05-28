package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/service"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"
	"github.com/spf13/viper"

	profileApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/profile"

	_ "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/docs"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Sad Islands API
// @version 1.0
// @description This is a super game.

// @host 104.248.28.45
// @BasePath /api/v1

func main() {
	clientId := flag.String("client_id", "", "client id")
	clientSecret := flag.String("client_secret", "", "client secret")
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
	port := viper.GetString("server.port")
	profileName := viper.GetString("profile.name")
	sessionName := viper.GetString("session.name")

	// Start watching for updates in consul
	resolver.RegisterDefault(consulAddr, consulPort, 3*time.Second)

	// Register prometheus metrics
	metric.RegisterAccessHitsMetric("api_service")

	// Create zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileConn := service.Connect(consulAddr, profileName)
	defer profileConn.Close()
	profileManager := profile.NewProfileClient(profileConn)

	sessionConn := service.Connect(consulAddr, sessionName)
	defer sessionConn.Close()
	sessionManager := session.NewSessionClient(sessionConn)

	profileApi := profileApi.NewApi(&profileManager, &sessionManager, logger, *clientId, *clientSecret)
	profileApi.Router.Handler("GET", "/swagger/:file", httpSwagger.WrapHandler)
	profileApi.Router.ServeFiles("/upload/*filepath", http.Dir("upload"))

	logger.Info("Server exit",
		zap.Error(http.ListenAndServe(":"+port, profileApi.Router)))
}
