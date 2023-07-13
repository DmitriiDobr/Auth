package handler

import (
	"auth/internal/auth/types"
	"auth/pkg/fiberErr"

	"github.com/gofiber/fiber/v2"
	"time"
)

func (h *Handler) register(c *fiber.Ctx) error {
	var user types.User
	err := c.BodyParser(&user)
	err = h.service.Register(c.Context(), user, string(c.Response().Header.Header()), string(c.Request().Body()))
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(&user)
	}
	//switch {
	//case errors.Is(err, types.ErrUserCreation):
	//	return c.Status(fiber.StatusInternalServerError).SendStatus(500)
	//
	//case errors.Is(err, types.ErrTimeRunnedOut):
	//	return c.Status(fiber.StatusUnauthorized).SendStatus(401)
	//
	//case errors.Is(err, types.ErrCantSendMessage):
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	//default:
	//	return err
	//}
	return nil

}

func (h *Handler) signin(c *fiber.Ctx) error {
	var creds types.Login
	if err := c.BodyParser(&creds); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	expirationDate, token, tokenString, err := h.service.Login(c.Context(), creds,
		string(c.Response().Header.Header()), string(c.Request().Body()))
	if err == nil {
		cookie := fiber.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			Expires:  expirationDate,
			HTTPOnly: true,
		}

		c.Cookie(&cookie)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"access_token": token})

	}

	return fiberErr.From(c, err)
}

func (h *Handler) refresh(c *fiber.Ctx) error {
	token := c.Cookies("access_token", "mistake")
	expirationalTime, tokenString, err := h.service.Refresh(token)
	if err == nil {
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

	return fiberErr.From(c, err)
}

func (h *Handler) logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:    "access_token",
		Expires: time.Now(),
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": fiber.StatusOK})
}
