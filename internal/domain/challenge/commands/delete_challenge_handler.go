package commands

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
	"context"
	"errors"
	"log/slog"
)

type DeleteChallengeHandler struct {
	cqrs.CommandHandler[DeleteChallengeCommand]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewDeleteChallengeHandler(log *slog.Logger, cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *DeleteChallengeHandler {
	return &DeleteChallengeHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (h *DeleteChallengeHandler) Handle(ctx context.Context, command cqrs.Command) (interface{}, error) {
	h.log.Info("DeleteChallengeHandler")
	deleteChallengeCommand, ok := command.(*DeleteChallengeCommand)
	if !ok {
		return nil, errors.New("invalid command")
	}

	err := h.repo.Delete(deleteChallengeCommand.ChallengeID)
	if err != nil {
		return nil, err
	}
	return "successful deleted", nil
}
