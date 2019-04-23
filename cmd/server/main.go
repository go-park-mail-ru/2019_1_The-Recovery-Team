package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/pkg/resolver"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/http/rest/profile/api"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/delivery/grpc/service/session"

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
	consulAddr := "consul"
	consulPort := 8500
	resolver.RegisterDefault(consulAddr, consulPort, 5*time.Second)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger creation error")
	}
	defer logger.Sync()

	profileConn, err := grpc.Dial("srv://consul/profile-service", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Can't connect to profile service:", err)
	}
	defer profileConn.Close()

	sessionConn, err := grpc.Dial("srv://consul/session-service", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatal("Authentication service connection refused:", err)
	}
	defer sessionConn.Close()

	sessionManager := session.NewSessionClient(sessionConn)
	profileManager := profile.NewProfileClient(profileConn)

	profileApi := api.NewApi(&profileManager, &sessionManager, logger)
	profileApi.Router.Handler("GET", "/swagger/:file", httpSwagger.WrapHandler)
	profileApi.Router.ServeFiles("/upload/*filepath", http.Dir("upload"))

	log.Print(http.ListenAndServe(":"+port, profileApi.Router))
}
