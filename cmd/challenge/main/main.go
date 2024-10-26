package main

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/delievery/http"
	"challenge-service/internal/domain/challenge/delievery/http/handlers"
	"challenge-service/internal/infrastructure/database/postgres"
	"challenge-service/internal/infrastructure/repository"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var dbClient *gorm.DB
	cfg := config.MustLoadConfig("config/config.yaml")
	log := setupLogger(cfg.Env)
	log.Info("Logger started successfully")
	pgConnect := postgres.NewPostgresConnect(cfg)
	client, err := pgConnect.Connect()
	if err != nil {
		panic(err)
	}
	dbClient = client.(*gorm.DB)

	defer func(pgConnect postgres.PostgresConnectable, i interface{}) {
		err := pgConnect.CloseConnection(i)
		if err != nil {
			panic(err)
		}
	}(pgConnect, dbClient)
	challengeRepo := repository.NewChallengeRepository(cfg, log, dbClient)
	challengeHandlers := handlers.NewChallengesHandlers(cfg, log, challengeRepo)
	httpServer := http.NewHTTPServer(cfg, log, challengeHandlers)
	httpServer.Run()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev, envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
