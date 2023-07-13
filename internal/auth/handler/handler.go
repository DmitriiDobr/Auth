package handler

import (
	"auth/internal/auth/types"
	"context"
	"encoding/json"
	jwtBus "github.com/LdDl/fiber-jwt/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type serviceAction interface {
	Login(ctx context.Context, creds types.Login, header, body string) (time.Time, *jwt.Token, string, error)
	Register(ctx context.Context, user types.User, header, body string) error
	Refresh(token string) (time.Time, string, error)
}

type Handler struct {
	service serviceAction
	jwtBus  *jwtBus.FiberJWTMiddleware
}

func NewHandler(service serviceAction) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandlers(app *fiber.App) {
	api := app.Group("/api/auth/v1")
	api.Post("/register", h.register)
	api.Post("/login", h.signin)
	api.Use(accessMiddleware)
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339Nano,
		TimeZone:   "Asia/Shanghai",
		Done: func(c *fiber.Ctx, logString []byte) {
			var mp map[string]any
			if c.Response().StatusCode() != fiber.StatusOK {
				json.Unmarshal(c.Response().Body(), &mp)
				logrus.Error(c.Request().Header.String(), mp)
			}
		},
	}))
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
