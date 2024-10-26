package queries

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
	"context"
	"errors"
	"log/slog"
)

type GetAllChallengesFromTeamQueryHandler struct {
	cqrs.QueryHandler[GetAllChallengesFromTeamQuery]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewGetAllChallengesFromTeamQueryHandler(log *slog.Logger, cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *GetAllChallengesFromTeamQueryHandler {
	return &GetAllChallengesFromTeamQueryHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (handler *GetAllChallengesFromTeamQueryHandler) Handle(ctx context.Context, query cqrs.Query) (interface{}, error) {
	handler.log.Info("GetAllChallengesFromTeamQueryHandler")
	getAllChallengesFromTeamQuery, ok := query.(*GetAllChallengesFromTeamQuery)
	if !ok {
		return nil, errors.New("invalid query type")
	}

	result, err := handler.repo.GetAllChallengesFromTeam(getAllChallengesFromTeamQuery.TeamID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
