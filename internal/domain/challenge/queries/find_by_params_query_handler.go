package queries

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
	"context"
	"errors"
	"log/slog"
)

type FindByParamsQueryHandler struct {
	cqrs.QueryHandler[FindByParamsQuery]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewFindByParamsQueryHandler(log *slog.Logger, cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *FindByParamsQueryHandler {
	return &FindByParamsQueryHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}

func (handler *FindByParamsQueryHandler) Handle(ctx context.Context, query cqrs.Query) (interface{}, error) {
	handler.log.Info("FindByParamsQueryHandler")
	findByParamsQuery, ok := query.(*FindByParamsQuery)
	if !ok {
		return nil, errors.New("invalid query type")
	}

	if findByParamsQuery.Params == nil {
		return nil, errors.New("missing parameters")
	}

	result, err := handler.repo.FindByParams(findByParamsQuery.Params)
	if err != nil {
		return nil, err
	}
	return result, nil
}
