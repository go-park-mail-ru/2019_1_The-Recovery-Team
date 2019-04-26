package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"

	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	gameApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"google.golang.org/grpc/balancer/roundrobin"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

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
	port := viper.GetInt("game.port")
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
	gameManager := usecase.NewGameInteractor(game.NewGameRepo(logger))

	go gameManager.Run()

	api := gameApi.NewApi(&profileManager, &sessionManager, gameManager, logger)

	log.Print(http.ListenAndServe(":"+strconv.Itoa(port), api.Router))
}
