package errors

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler(logger *slog.Logger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		var appError *Error
		var fiberError *fiber.Error

		logger.Error("Error", "error", err)

		switch true {
		case errors.As(err, &appError):
			break
		case errors.As(err, &fiberError):
			var errorCode ErrorCode
			switch fiberError.Code {
			case fiber.StatusNotFound:
				errorCode = ErrorCodeNotFound
			default:
				errorCode = ErrorCodeServerError
			}
			appError = NewHTTPError(fiberError.Code, errorCode, fiberError.Message)
		default:
			appError = NewHTTPError(fiber.StatusInternalServerError, ErrorCodeServerError, err.Error())
		}

		return ctx.Status(appError.Status).JSON(&appError)
	}
}
