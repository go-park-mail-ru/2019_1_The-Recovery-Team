package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"

	profileApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
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

func main() {
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
	resolver.RegisterDefault(consulAddr, consulPort, 5*time.Second)

	// Register prometheus metrics
	metric.RegisterAccessHitsMetric("api_service")
	prometheus.MustRegister(metric.AccessHits)

	port := viper.GetString("server.port")
	profileName := viper.Get("profile.name")
	sessionName := viper.Get("session.name")

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileConn, err := grpc.Dial(fmt.Sprintf("srv://%s/%s", consulAddr, profileName),
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name), grpc.WithBlock())
	if err != nil {
		log.Fatal("Can't connect to profile service:", err)
	}
	defer profileConn.Close()

	sessionConn, err := grpc.Dial(fmt.Sprintf("srv://%s/%s", consulAddr, sessionName),
		grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name), grpc.WithBlock())
	if err != nil {
		log.Fatal("Authentication service connection refused:", err)
	}
	defer sessionConn.Close()

	sessionManager := session.NewSessionClient(sessionConn)
	profileManager := profile.NewProfileClient(profileConn)

	profileApi := profileApi.NewApi(&profileManager, &sessionManager, logger)
	profileApi.Router.Handler("GET", "/swagger/:file", httpSwagger.WrapHandler)
	profileApi.Router.ServeFiles("/upload/*filepath", http.Dir("upload"))

	logger.Info("Server exit",
		zap.Error(http.ListenAndServe(":"+port, profileApi.Router)))
}
