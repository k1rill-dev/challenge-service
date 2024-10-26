package commands

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/entity"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
	"context"
	"errors"
	"log/slog"
)

type UpdateChallengeHandler struct {
	cqrs.CommandHandler[UpdateChallengeCommand]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewUpdateChallengeHandler(log *slog.Logger, cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *UpdateChallengeHandler {
	return &UpdateChallengeHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (h *UpdateChallengeHandler) Handle(ctx context.Context, command cqrs.Command) (interface{}, error) {
	h.log.Info("UpdateChallengeHandler")
	updateChallengeCommand, ok := command.(*UpdateChallengeCommand)
	if !ok {
		return nil, errors.New("invalid command")
	}
	challenge := entity.AuthenticationChallenge{
		ID: updateChallengeCommand.ChallengeID,
	}

	if updateChallengeCommand.Name != nil {
		challenge.Name = *updateChallengeCommand.Name
	}
	if updateChallengeCommand.Icon != nil {
		challenge.Icon = *updateChallengeCommand.Icon
	}
	if updateChallengeCommand.Image != nil {
		challenge.Image = *updateChallengeCommand.Image
	}
	if updateChallengeCommand.Description != nil {
		challenge.Description = *updateChallengeCommand.Description
	}
	if updateChallengeCommand.StartDate != nil {
		challenge.StartDate = *updateChallengeCommand.StartDate
	}
	if updateChallengeCommand.EndDate != nil {
		challenge.EndDate = *updateChallengeCommand.EndDate
	}
	if updateChallengeCommand.Type != nil {
		challenge.Type = *updateChallengeCommand.Type
	}
	if updateChallengeCommand.IsTeam != nil {
		challenge.IsTeam = *updateChallengeCommand.IsTeam
	}
	if updateChallengeCommand.CreatorID != nil {
		challenge.CreatorID = *updateChallengeCommand.CreatorID
	}

	result, err := h.repo.Update(challenge)
	if err != nil {
		return nil, err
	}
	return result, nil
}
