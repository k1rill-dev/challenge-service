package handlers

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/commands"
	"challenge-service/internal/domain/challenge/entity"
	"challenge-service/internal/domain/challenge/queries"
	"challenge-service/internal/infrastructure/lib/fabric"
	"challenge-service/internal/infrastructure/lib/log"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
)

type ChallengesHandlers struct {
	cfg           *config.Config
	log           *slog.Logger
	handlerFabric *fabric.HandlerFabric
}

func NewChallengesHandlers(cfg *config.Config, log *slog.Logger, handlerFabric *fabric.HandlerFabric) *ChallengesHandlers {
	return &ChallengesHandlers{
		cfg:           cfg,
		log:           log,
		handlerFabric: handlerFabric,
	}
}

func (h *ChallengesHandlers) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *ChallengesHandlers) CreateChallenge(c *gin.Context) {
	var challenge entity.AuthenticationChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		h.log.Error("Error binding JSON:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	randomID := rand.Int64()
	command := commands.NewCreateChallengeCommand(randomID, &challenge.Name, &challenge.Icon, &challenge.Description,
		&challenge.Interests, &challenge.EndDate, &challenge.Type, &challenge.IsTeam, &challenge.CreatorID)

	handler, err := h.handlerFabric.GetCommandHandler(command)
	if err != nil {
		h.log.Error("Error getting command handler:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := handler.Handle(context.Background(), command)
	if err != nil {
		h.log.Error("Error handling command:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Получение всех вызовов
func (h *ChallengesHandlers) GetAllChallenges(c *gin.Context) {
	query := queries.NewFindAllQuery()
	handler, err := h.handlerFabric.GetCommandHandler(query)
	if err != nil {
		h.log.Error("Error getting query handler:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := handler.Handle(context.Background(), query)
	if err != nil {
		h.log.Error("Error handling query:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Обновление вызова
func (h *ChallengesHandlers) UpdateChallenge(c *gin.Context) {
	var updateCommand commands.UpdateChallengeCommand
	if err := c.ShouldBindJSON(&updateCommand); err != nil {
		h.log.Error("Error binding JSON:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	handler, err := h.handlerFabric.GetCommandHandler(&updateCommand)
	if err != nil {
		h.log.Error("Error getting command handler:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := handler.Handle(context.Background(), &updateCommand)
	if err != nil {
		h.log.Error("Error handling command:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Удаление вызова
func (h *ChallengesHandlers) DeleteChallenge(c *gin.Context) {
	idParam := c.Param("id")
	challengeID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		h.log.Error("Error parsing challenge ID:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge ID"})
		return
	}

	command := commands.NewDeleteChallengeCommand(rand.Int64(), challengeID)
	handler, err := h.handlerFabric.GetCommandHandler(command)
	if err != nil {
		h.log.Error("Error getting command handler:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = handler.Handle(context.Background(), command)
	if err != nil {
		h.log.Error("Error handling command:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Получение вызовов пользователя
func (h *ChallengesHandlers) GetAllChallengesFromUser(c *gin.Context) {
	userID := c.Param("user_id")
	query := queries.NewGetAllChallengesFromUserQuery(rand.Int64(), userID)
	handler, err := h.handlerFabric.GetCommandHandler(query)
	if err != nil {
		h.log.Error("Error getting query handler:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := handler.Handle(context.Background(), query)
	if err != nil {
		h.log.Error("Error handling query:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// Получение вызовов команды
func (h *ChallengesHandlers) GetAllChallengesFromTeam(c *gin.Context) {
	teamID := c.Param("team_id")
	query := queries.NewGetAllChallengesFromTeamQuery(rand.Int64(), teamID)
	handler, err := h.handlerFabric.GetCommandHandler(query)
	if err != nil {
		h.log.Error("Error getting query handler:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := handler.Handle(context.Background(), query)
	if err != nil {
		h.log.Error("Error handling query:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
