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

type CreateChallengeHandler struct {
	cqrs.CommandHandler[CreateChallengeCommand]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewCreateChallengeHandler(log *slog.Logger, cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *CreateChallengeHandler {
	return &CreateChallengeHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (c *CreateChallengeHandler) Handle(ctx context.Context, command cqrs.Command) (interface{}, error) {
	c.log.Info("CreateChallengeHandler")
	createChallengeCommand, ok := command.(*CreateChallengeCommand)
	if !ok {
		return nil, errors.New("invalid command")
	}
	challenge := entity.AuthenticationChallenge{
		ID:          createChallengeCommand.AggregateID,
		Name:        createChallengeCommand.Name,
		Icon:        createChallengeCommand.Icon,
		Image:       createChallengeCommand.Image,
		Description: createChallengeCommand.Description,
		Interests:   createChallengeCommand.Interests,
		EndDate:     createChallengeCommand.EndDate,
		Type:        createChallengeCommand.Type,
		IsTeam:      createChallengeCommand.IsTeam,
		CreatorID:   createChallengeCommand.CreatorID,
	}
	result, err := c.repo.Create(challenge)
	if err != nil {
		return nil, err
	}
	return result, nil
}
