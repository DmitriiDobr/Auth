package handler

import (
	"auth/internal/auth/service"
	"auth/internal/auth/types"
	"errors"
	"fmt"
	kafkaNotification "github.com/DmitriiDobr/kafkaNotification/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func (h *Handler) register(c *fiber.Ctx) error {
	var user types.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&user)
	}
	_, err = h.service.CreateUser(c.Context(), user)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	}
	msg := kafkaNotification.Message{
		UserID: user.Id,
		Status: kafkaNotification.Success,
		Header: string(c.Response().Header.Header()),
		Body:   string(c.Request().Body()),
	}

	err = h.kafkaClient.Notify(c.Context(), msg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": kafkaNotification.Error})
	}
	return c.Status(fiber.StatusOK).JSON(&user)

}

func (h *Handler) signin(c *fiber.Ctx) error {
	var creds types.Login
	if err := c.BodyParser(&creds); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	user, err := h.service.GetUser(c.Context(), creds.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if user.Password != service.GeneratePasswordHash(creds.Password, h.salt) {
		return fiber.NewError(fiber.StatusBadRequest, "Ваш пароль некорректен!")
	}
	expirationDate := time.Now().Add(30 * time.Second)
	claims := &types.Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(token)
	tokenString, err := token.SignedString(h.jwtKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при подписании токена!")
	}
	fmt.Println(tokenString)
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Expires:  expirationDate,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	msg := kafkaNotification.Message{
		UserID: user.Id,
		Status: kafkaNotification.Success,
		Header: string(c.Response().Header.Header()),
		Body:   string(c.Request().Body()),
	}

	err = h.kafkaClient.Notify(c.Context(), msg)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": token})
}

func (h *Handler) refresh(c *fiber.Ctx) error {
	token := c.Cookies("access_token", "mistake")
	if token == "mistake" {
		return fiber.NewError(fiber.StatusUnauthorized, "Пользователь не авторизован!")
	}
	claims := &types.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return h.jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return fiber.NewError(fiber.StatusUnauthorized, "Пользователь не авторизован!")
		}
		return fiber.NewError(fiber.StatusBadRequest)
	}
	if !tkn.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Пользователь не авторизован!")
	}
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return fiber.NewError(fiber.StatusBadRequest, "Время токена просрочено")

	}
	expirationalTime := time.Now().Add(30 * time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(expirationalTime)
	tokenUpdate := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenUpdate.SignedString(h.jwtKey)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при подписании токена")
	}
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Expires:  expirationalTime,
		HTTPOnly: true,
		Secure:   true,
	}

	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fiber.StatusOK})
}

func (h *Handler) logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:    "access_token",
		Expires: time.Now(),
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fiber.StatusOK})
}
