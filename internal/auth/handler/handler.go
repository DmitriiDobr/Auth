package handler

import (
	"auth/internal/auth/types"
	"context"
	kafkaNotification "github.com/DmitriiDobr/kafkaNotification/pkg"
	jwt "github.com/LdDl/fiber-jwt/v2"
	"github.com/gofiber/fiber/v2"
)

type Iservice interface {
	CreateUser(ctx context.Context, user types.User) (int, error)
	GetUser(ctx context.Context, username string) (*types.User, error)
}

type Handler struct {
	service     Iservice
	salt        string
	jwtKey      string
	jwtBus      *jwt.FiberJWTMiddleware
	kafkaClient *kafkaNotification.Client
}

func NewHandler(service Iservice, kafkaClient *kafkaNotification.Client, salt, jwtKey string) *Handler {
	return &Handler{service: service, kafkaClient: kafkaClient, salt: salt, jwtKey: jwtKey}
}

func (h *Handler) RegisterHandlers(app *fiber.App) {
	api := app.Group("/api/auth/v1")
	api.Post("/register", h.register)
	api.Post("/login", h.signin)
	api.Use(accessMiddleware)
	api.Get("/refresh", h.refresh)
	api.Get("/logout", h.logout)

}

func accessMiddleware(c *fiber.Ctx) error {
	accessToken := c.Cookies("access_token", "unauthorized")
	//fmt.Println(accessToken)
	if accessToken == "unauthorized" {
		// Посредник прерывает цепочку обработки запроса.
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Пользователь имеет доступ, продолжаем выполнение запроса.
	return c.Next()
}

//fiber проверяет есть ли
