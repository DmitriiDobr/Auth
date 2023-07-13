package fiberErr

import (
	"auth/pkg/errors"
	"github.com/gofiber/fiber/v2"
)

func From(c *fiber.Ctx, err error) error {
	customErr, ok := err.(*errors.CustomError)
	if !ok {
		c.Status(fiber.StatusInternalServerError)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	c.Status(customErr.StatusCode)
	return c.JSON(map[string]any{"code": customErr.StatusCode, "message": customErr.Error()})
}
