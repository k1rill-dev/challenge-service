package http

import (
	"challenge-service/config"
	"challenge-service/docs"
	"challenge-service/internal/domain/challenge/delievery/http/handlers"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
	"github.com/swaggo/files"       // swagger embedded files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
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
	router.GET("/pingpong", h.challengesHandlers.Ping)

	// Группа маршрутов, защищенных AuthMiddleware
	api := router.Group("/api")
	api.Use(AuthMiddleware(h.cfg))

	// Роуты для управления вызовами (challenges)
	challenges := api.Group("/challenges")
	{
		// Создание нового вызова
		challenges.POST("/", h.challengesHandlers.CreateChallenge)

		// Получение всех вызовов
		challenges.GET("/", h.challengesHandlers.GetAllChallenges)

		// Обновление вызова
		challenges.PUT("/:id", h.challengesHandlers.UpdateChallenge)

		// Удаление вызова
		challenges.DELETE("/:id", h.challengesHandlers.DeleteChallenge)

		// Получение вызовов пользователя
		challenges.GET("/user/:user_id", h.challengesHandlers.GetAllChallengesFromUser)

		// Получение вызовов команды
		challenges.GET("/team/:team_id", h.challengesHandlers.GetAllChallengesFromTeam)
	}
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Запуск сервера
	err := router.Run(":8000")
	if err != nil {
		h.log.Error("Failed to run server:", err)
		panic(err)
	}
	//"D:\GoProjects\challenge-service\internal\domain\challenge\delievery\http\handlers"
}
