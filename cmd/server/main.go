package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"

	profileApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc"

	_ "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Sad Islands API
// @version 1.0
// @description This is a super game.

// @host 104.248.28.45
// @BasePath /api/v1

func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("build/config/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config files:", err)
	}

	addr := viper.GetString("consul.address")
	port := viper.GetInt("consul.port")
	resolver.RegisterDefault(addr, port, 5*time.Second)
}

func main() {
	port := viper.GetString("server.port")
	profileName := viper.Get("profile.name")
	sessionName := viper.Get("session.name")
	consulAddr := viper.Get("consul.address")

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileConn, err := grpc.Dial(fmt.Sprintf("srv://%s/%s", consulAddr, profileName),
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Can't connect to profile service:", err)
	}
	defer profileConn.Close()

	sessionConn, err := grpc.Dial(fmt.Sprintf("srv://%s/%s", consulAddr, sessionName),
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Authentication service connection refused:", err)
	}
	defer sessionConn.Close()

	sessionManager := session.NewSessionClient(sessionConn)
	profileManager := profile.NewProfileClient(profileConn)

	profileApi := profileApi.NewApi(&profileManager, &sessionManager, logger)
	profileApi.Router.Handler("GET", "/swagger/:file", httpSwagger.WrapHandler)
	profileApi.Router.ServeFiles("/upload/*filepath", http.Dir("upload"))

	log.Print(http.ListenAndServe(":"+port, profileApi.Router))
}
