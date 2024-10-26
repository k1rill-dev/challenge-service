package main

import (
	"challenge-service/config"
	"challenge-service/internal/infrastructure/database/postgres"
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
	cfg := config.MustLoadConfig("config/config.yaml")
	log := setupLogger(cfg.Env)
	log.Info("Logger started successfully")
	pgConnect := postgres.NewPostgresConnect(cfg)
	dbClient, err := pgConnect.Connect()
	if err != nil {
		panic(err)
	}
	dbClient = dbClient.(*gorm.DB)

	defer func(pgConnect postgres.PostgresConnectable, i interface{}) {
		err := pgConnect.CloseConnection(i)
		if err != nil {
			panic(err)
		}
	}(pgConnect, dbClient)

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
