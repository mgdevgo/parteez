package errors

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler(logger *slog.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		logger.Error("Error", "error", err)

		var pzError *Error
		if errors.As(err, &pzError) {
			return ctx.Status(pzError.Status).JSON(&pzError)
		}

		var fiberError *fiber.Error
		if errors.As(err, &fiberError) {
			return ctx.Status(fiberError.Code).JSON(NewHTTPError(fiberError.Code, "", fiberError.Message))
		}

		// Default internal server error
		result := NewHTTPError(fiber.StatusInternalServerError, "", err.Error())

		return ctx.Status(result.Status).JSON(result)
	}
}
