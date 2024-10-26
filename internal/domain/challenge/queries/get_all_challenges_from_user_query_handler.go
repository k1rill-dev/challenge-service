package queries

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
	"context"
	"errors"
	"log/slog"
)

type GetAllChallengesFromUserQueryHandler struct {
	cqrs.QueryHandler[GetAllChallengesFromUserQuery]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewGetAllChallengesFromUserQueryHandler(log *slog.Logger, cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *GetAllChallengesFromUserQueryHandler {
	return &GetAllChallengesFromUserQueryHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (handler *GetAllChallengesFromUserQueryHandler) Handle(ctx context.Context, query cqrs.Query) (interface{}, error) {
	handler.log.Info("GetAllChallengesFromUserQueryHandler")
	getAllChallengesFromUserQuery, ok := query.(*GetAllChallengesFromUserQuery)
	if !ok {
		return nil, errors.New("invalid query type")
	}

	result, err := handler.repo.GetAllChallengesFromUser(getAllChallengesFromUserQuery.UserID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
