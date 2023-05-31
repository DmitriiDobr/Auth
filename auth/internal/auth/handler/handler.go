package handler

import (
	"auth/internal/auth/service"
	"fmt"
	jwt "github.com/LdDl/fiber-jwt/v2"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *service.AuthService
	jwtBus  *jwt.FiberJWTMiddleware
}

func NewHandler(service *service.AuthService) *Handler {
	return &Handler{service: service}
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
	fmt.Println(accessToken)
	if accessToken == "unauthorized" {
		// Посредник прерывает цепочку обработки запроса.
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Пользователь имеет доступ, продолжаем выполнение запроса.
	return c.Next()
}

//fiber проверяет есть ли
