package handlers

import (
	"challenge-service/config"
	repo "challenge-service/internal/domain/challenge/usecases/repository_interface"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type ChallengesHandlers struct {
	cfg           *config.Config
	log           *slog.Logger
	challengeRepo repo.ChallengeRepositoryInterface
}

func NewChallengesHandlers(cfg *config.Config, log *slog.Logger, challengeRepo repo.ChallengeRepositoryInterface) *ChallengesHandlers {
	return &ChallengesHandlers{
		cfg:           cfg,
		log:           log,
		challengeRepo: challengeRepo,
	}
}

func (h *ChallengesHandlers) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
