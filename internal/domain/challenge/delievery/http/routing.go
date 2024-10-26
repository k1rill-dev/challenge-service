package http

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/delievery/http/handlers"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type HTTPServer struct {
	cfg                *config.Config
	log                *slog.Logger
	challengesHandlers *handlers.ChallengesHandlers
}

func NewHTTPServer(cfg *config.Config, log *slog.Logger, challengeHandlers *handlers.ChallengesHandlers) *HTTPServer {
	return &HTTPServer{
		cfg:                cfg,
		log:                log,
		challengesHandlers: challengeHandlers,
	}
}

func (h *HTTPServer) Run() {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.GET("/pingpong", h.challengesHandlers.Ping)

	//router.POST("/lobbies", h.challengesHandlers.CreateLobby)
	//router.GET("/lobbies/:lobbyId", h.lobbyHandlers.GetLobby)
	//router.GET("/users/:userId/lobbies", h.lobbyHandlers.GetAllLobbiesFromUser)
	//router.POST("/lobbies/:lobbyId/users", h.lobbyHandlers.AddUsersToLobby)
	//router.DELETE("/lobbies/:lobbyId", h.lobbyHandlers.DeleteLobby)
	//router.DELETE("/lobbies/:lobbyId/users/:userId", h.lobbyHandlers.RemoveUserFromLobby)
	//router.PUT("/lobbies/:lobbyId", h.lobbyHandlers.UpdateLobby)
	//router.GET("/lobbies/:lobbyId/users", h.lobbyHandlers.GetAllUsersFromLobby)

	err := router.Run(":8000")
	if err != nil {
		panic(err)
	}
}
