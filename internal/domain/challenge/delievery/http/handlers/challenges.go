// Package handlers ...
// @title Challenge Service API
// @version 1.0
// @description This is a sample server for managing challenges.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

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

type PingResponse struct {
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
}

// Ping
// @Summary      Check service health
// @Description  Responds with a "pong" message to check service availability
// @Tags         Health
// @Produce      json
// @Success      200  {object} PingResponse
// @Router       /pingpong [get]
func (h *ChallengesHandlers) Ping(c *gin.Context) {
	c.JSON(200, PingResponse{
		Message: "pong",
	})
}

// CreateChallenge
// @Summary      Create a new challenge
// @Description  Creates a new challenge with the provided data
// @Tags         Challenges
// @Accept       json
// @Produce      json
// @Param        challenge  body  entity.AuthenticationChallenge  true  "Challenge Data"
// @Success      201  {object}  entity.AuthenticationChallenge // Изменен код успешного ответа
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges [post]
func (h *ChallengesHandlers) CreateChallenge(c *gin.Context) {
	var challenge entity.AuthenticationChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		h.log.Error("Error binding JSON:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	randomID := rand.Int64()
	command := commands.NewCreateChallengeCommand(randomID, &challenge.Name, &challenge.Icon, &challenge.Description,
		&challenge.EndDate, &challenge.Type, &challenge.IsTeam, &challenge.CreatorID)

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
	c.JSON(http.StatusCreated, result) // Изменен код успешного ответа
}

// GetAllChallenges
// @Summary      Retrieve all challenges
// @Description  Fetches a list of all challenges
// @Tags         Challenges
// @Produce      json
// @Success      200  {array}  entity.AuthenticationChallenge
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges [get]
func (h *ChallengesHandlers) GetAllChallenges(c *gin.Context) {
	query := queries.NewFindAllQuery()
	handler, err := h.handlerFabric.GetQueryHandler(query)
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

// UpdateChallenge
// @Summary      Update an existing challenge
// @Description  Updates the details of an existing challenge
// @Tags         Challenges
// @Accept       json
// @Produce      json
// @Param        id         path     int64                      true  "Challenge ID"
// @Param        challenge  body     entity.AuthenticationChallenge true  "Updated Challenge Data"
// @Success      200  {object}  interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/{id} [put]
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

type DeleteChallengeResponse struct {
	Message string `json:"message"`
}

// DeleteChallenge
// @Summary      Delete a challenge
// @Description  Removes a challenge by its ID
// @Tags         Challenges
// @Param        id   path     int64  true  "Challenge ID"
// @Success      200  {object}  DeleteChallengeResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/{id} [delete]
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

// GetAllChallengesFromUser
// @Summary      Get challenges for a user
// @Description  Retrieves all challenges associated with a specific user
// @Tags         Challenges
// @Param        user_id  path     string  true  "User ID"
// @Produce      json
// @Success      200  {array}  entity.AuthenticationChallenge
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/user/{user_id} [get]
func (h *ChallengesHandlers) GetAllChallengesFromUser(c *gin.Context) {
	userID := c.Param("user_id")
	query := queries.NewGetAllChallengesFromUserQuery(rand.Int64(), userID)
	handler, err := h.handlerFabric.GetQueryHandler(query)
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

// GetAllChallengesFromTeam
// @Summary      Get challenges for a team
// @Description  Retrieves all challenges associated with a specific team
// @Tags         Challenges
// @Param        team_id  path     string  true  "Team ID"
// @Produce      json
// @Success      200  {array}  entity.AuthenticationChallenge
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/team/{team_id} [get]
func (h *ChallengesHandlers) GetAllChallengesFromTeam(c *gin.Context) {
	teamID := c.Param("team_id")
	query := queries.NewGetAllChallengesFromTeamQuery(rand.Int64(), teamID)
	handler, err := h.handlerFabric.GetQueryHandler(query)
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
