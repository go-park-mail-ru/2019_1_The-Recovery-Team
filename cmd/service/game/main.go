package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/service"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/metric"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/pkg/resolver"
	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	gameApi "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/api/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/memory/game"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"

	"go.uber.org/zap"
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
	port := viper.GetInt("game.port")
	profileName := viper.GetString("profile.name")
	sessionName := viper.GetString("session.name")

	// Start watching for updates in consul
	resolver.RegisterDefault(consulAddr, consulPort, 5*time.Second)

	// Register prometheus metrics
	metric.RegisterTotalRoomsMetric("game_service")
	metric.RegisterTotalPlayersMetric("game_service")

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

	gameManager := usecase.NewGameInteractor(game.NewGameRepo(logger))

	api := gameApi.NewApi(&profileManager, &sessionManager, gameManager, logger)

	logger.Info("Game exit",
		zap.Error(http.ListenAndServe(":"+strconv.Itoa(port), api.Router)))
}
