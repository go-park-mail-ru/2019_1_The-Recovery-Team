package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"
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
	metric.RegisterTotalRoomsMetric("game_service")
	metric.RegisterTotalPlayersMetric("game_service")

	port := viper.GetInt("game.port")
	profileName := viper.Get("profile.name")
	sessionName := viper.Get("session.name")

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

	api := gameApi.NewApi(&profileManager, &sessionManager, gameManager, logger)

	logger.Info("Game exit",
		zap.Error(http.ListenAndServe(":"+strconv.Itoa(port), api.Router)))
}
