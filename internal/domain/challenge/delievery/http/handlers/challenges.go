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
	"challenge-service/internal/domain/challenge/usecases/repository_interface"
	"challenge-service/internal/infrastructure/lib/fabric"
	"challenge-service/internal/infrastructure/lib/log"
	"challenge-service/internal/infrastructure/lib/save_photo"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
)

type ChallengesHandlers struct {
	cfg           *config.Config
	log           *slog.Logger
	handlerFabric *fabric.HandlerFabric
	repo          repository_interface.ChallengeRepositoryInterface
}

func NewChallengesHandlers(cfg *config.Config, log *slog.Logger, handlerFabric *fabric.HandlerFabric,
	repo repository_interface.ChallengeRepositoryInterface) *ChallengesHandlers {
	return &ChallengesHandlers{
		cfg:           cfg,
		log:           log,
		handlerFabric: handlerFabric,
		repo:          repo,
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Summary      Create a new challenge
// @Description  Creates a new challenge with the provided data
// @Tags         Challenges
// @Accept       multipart/form-data
// @Produce      json
// @Param        challenge  body  entity.AuthenticationChallenge  true  "Challenge Data"
// @Param        image      formData  file  true  "Image File"
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
	image, header, err := c.Request.FormFile("image")
	if err != nil {
		h.log.Error("Error retrieving image:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image not provided"})
		return
	}
	defer image.Close()
	icon, headerIcon, err := c.Request.FormFile("icon")
	if err != nil {
		h.log.Error("Error retrieving icon:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Icon not provided"})
	}
	defer icon.Close()
	imageBytes, err := io.ReadAll(image)
	if err != nil {
		h.log.Error("Error reading file:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
		return
	}
	iconBytes, err := io.ReadAll(icon)
	if err != nil {
		h.log.Error("Error reading file:", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read icon"})
		return
	}
	s3Client := save_photo.NewS3Client(h.cfg, h.log)
	urlImage, err := s3Client.UploadFile(imageBytes, header.Filename)
	if err != nil {
		h.log.Error("error while saving photo:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	urlIcon, err := s3Client.UploadFile(iconBytes, headerIcon.Filename)
	challenge.Image = urlImage
	challenge.Icon = urlIcon
	if err != nil {
		h.log.Error("error while saving icon:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Summary      Update an existing challenge
// @Description  Updates the details of an existing challenge
// @Tags         Challenges
// @Accept       multipart/form-data
// @Produce      json
// @Param        id         path      int64                      true  "Challenge ID"
// @Param        challenge  body      entity.AuthenticationChallenge false "Updated Challenge Data"
// @Param        image      formData  file  false  "New Image File"
// @Param        icon       formData  file  false  "New Icon File"
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

	s3Client := save_photo.NewS3Client(h.cfg, h.log)

	// Проверка и загрузка нового изображения, если оно предоставлено
	if image, header, err := c.Request.FormFile("image"); err == nil {
		defer image.Close()
		imageBytes, err := io.ReadAll(image)
		if err != nil {
			h.log.Error("Error reading image file:", log.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
			return
		}
		urlImage, err := s3Client.UploadFile(imageBytes, header.Filename)
		if err != nil {
			h.log.Error("Error uploading image to S3:", log.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updateCommand.Image = &urlImage
	}

	// Проверка и загрузка новой иконки, если она предоставлена
	if icon, headerIcon, err := c.Request.FormFile("icon"); err == nil {
		defer icon.Close()
		iconBytes, err := io.ReadAll(icon)
		if err != nil {
			h.log.Error("Error reading icon file:", log.Err(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read icon"})
			return
		}
		urlIcon, err := s3Client.UploadFile(iconBytes, headerIcon.Filename)
		if err != nil {
			h.log.Error("Error uploading icon to S3:", log.Err(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updateCommand.Icon = &urlIcon
	}

	// Обработка команды обновления
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

// RegisterUser
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Summary      Register user on challenge
// @Description  Register user on challenge
// @Tags         Challenges
// @Produce      json
// @Success      200  {array}  entity.AuthenticationParticipant
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/user/register [post]
func (h *ChallengesHandlers) RegisterUser(c *gin.Context) {
	userID, ok := c.Get("user_id")
	var challenge entity.AuthenticationChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		h.log.Error("Error binding JSON:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		h.log.Error("not auth", log.Err(errors.New("not authorized")))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	id := userID.(int64)
	response, err := h.repo.RegisterUserOnChallenge(id, challenge)
	if err != nil {
		h.log.Error("error while register on challenge", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while register on challenge"})
		return
	}
	c.JSON(http.StatusOK, response)
}

// RegisterTeam
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Summary      Register team on challenge
// @Description  Register team on challenge
// @Tags         Challenges
// @Param        team_id  path     string  true  "Team ID"
// @Produce      json
// @Success      200  {array}  entity.AuthenticationParticipant
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/team/register/{team_id} [post]
func (h *ChallengesHandlers) RegisterTeam(c *gin.Context) {
	_, ok := c.Get("user_id")
	if !ok {
		h.log.Error("not auth", log.Err(errors.New("not authorized")))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
	teamID := c.Param("team_id")
	var challenge entity.AuthenticationChallenge
	if err := c.ShouldBindJSON(&challenge); err != nil {
		h.log.Error("Error binding JSON:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := strconv.ParseInt(teamID, 10, 64)
	if err != nil {
		h.log.Error("Error parsing team ID:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := h.repo.RegisterTeamOnChallenge(id, challenge)
	if err != nil {
		h.log.Error("error while register on challenge", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while register on challenge"})
		return
	}
	c.JSON(http.StatusOK, response)
}

// CloseChallenge
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Summary      Close challenge
// @Description  This method closes challenge and send message to winner
// @Tags         Challenges
// @Param        challenge_id  path     string  true  "Challenge ID"
// @Produce      json
// @Success      200  {array}  entity.AuthenticationChallenge
// @Failure      500  {object}  ErrorResponse
// @Router       /challenges/close/{challenge_id} [post]
func (h *ChallengesHandlers) CloseChallenge(c *gin.Context) {
	challengeID, err := strconv.ParseInt(c.Param("challenge_id"), 10, 64)
	if err != nil {
		h.log.Error("Error parsing challenge ID:", log.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	response, err := h.repo.CloseChallenge(challengeID)
	if err != nil {
		h.log.Error("error while closing challenge", log.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error while closing challenge"})
	}
	c.JSON(http.StatusOK, response)
}
