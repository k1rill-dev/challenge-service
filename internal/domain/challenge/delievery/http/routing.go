package http

import (
	"challenge-service/config"
	"challenge-service/internal/domain/challenge/delievery/http/handlers"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
	"strings"
	"time"
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

// AuthMiddleware - middleware для проверки авторизации
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получить токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		// Проверить токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return cfg.SecretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Проверка истечения срока действия токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if exp, ok := claims["exp"].(float64); ok && time.Unix(int64(exp), 0).Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
			c.Set("userID", claims["user_id"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (h *HTTPServer) Run() {
	router := gin.Default()

	router.Use(gin.Recovery())
	myGroup := router.Group("/")
	myGroup.Use(AuthMiddleware(h.cfg))
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
