package main

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/commands"
	"challenge-service/internal/domain/challenge/delievery/http"
	"challenge-service/internal/domain/challenge/delievery/http/handlers"
	"challenge-service/internal/domain/challenge/queries"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/database/postgres"
	"challenge-service/internal/infrastructure/lib/fabric"
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
	handlerFabric := fabric.NewHandlerFabric()
	initializeHandlers(handlerFabric, log, cfg, challengeRepo)
	challengeHandlers := handlers.NewChallengesHandlers(cfg, log, handlerFabric, challengeRepo)
	httpServer := http.NewHTTPServer(cfg, log, challengeHandlers)
	httpServer.Run()
}

func initializeHandlers(
	handlerFabric *fabric.HandlerFabric,
	log *slog.Logger,
	config *config.Config,
	companyRepo repository_interface.ChallengeRepositoryInterface) {
	createChallengeHandler := commands.NewCreateChallengeHandler(log, config, companyRepo)
	updateChallengeHandler := commands.NewUpdateChallengeHandler(log, config, companyRepo)
	deleteChallengeHandler := commands.NewDeleteChallengeHandler(log, config, companyRepo)
	findAllHandler := queries.NewFindAllQueryHandler(log, config, companyRepo)
	findByParamsHandler := queries.NewFindByParamsQueryHandler(log, config, companyRepo)
	getAllChallengesFromTeamHandler := queries.NewGetAllChallengesFromTeamQueryHandler(log, config, companyRepo)
	getAllChallengesFromUserHandler := queries.NewGetAllChallengesFromUserQueryHandler(log, config, companyRepo)

	handlerFabric.RegisterCommandHandler(commands.NewEmptyCreateChallengeCommand(), createChallengeHandler)
	handlerFabric.RegisterCommandHandler(commands.NewEmptyUpdateChallengeCommand(), updateChallengeHandler)
	handlerFabric.RegisterCommandHandler(commands.NewEmptyDeleteChallengeCommand(), deleteChallengeHandler)
	handlerFabric.RegisterQueryHandler(queries.NewFindAllQuery(), findAllHandler)
	handlerFabric.RegisterQueryHandler(queries.NewEmptyFindByParamsQuery(), findByParamsHandler)
	handlerFabric.RegisterQueryHandler(queries.NewEmptyGetAllChallengesFromTeamQuery(), getAllChallengesFromTeamHandler)
	handlerFabric.RegisterQueryHandler(queries.NewEmptyGetAllChallengesFromUserQuery(), getAllChallengesFromUserHandler)

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
