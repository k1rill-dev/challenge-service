package queries

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/cqrs"
	"context"
	"log/slog"
)

type FindAllQueryHandler struct {
	cqrs.QueryHandler[FindAllQuery]
	log  *slog.Logger
	cfg  *config.Config
	repo repository_interface.ChallengeRepositoryInterface
}

func NewFindAllQueryHandler(log *slog.Logger,
	cfg *config.Config,
	repo repository_interface.ChallengeRepositoryInterface) *FindAllQueryHandler {
	return &FindAllQueryHandler{
		log:  log,
		cfg:  cfg,
		repo: repo,
	}
}
func (handler *FindAllQueryHandler) Handle(ctx context.Context, query cqrs.Query) (interface{}, error) {
	handler.log.Info("FindAllQueryHandler")
	_ = query.(*FindAllQuery)
	result, err := handler.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return result, nil
}
