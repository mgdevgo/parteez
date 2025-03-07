package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct {
	logger *slog.Logger
}

func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (h *ErrorHandler) Handle(ctx *fiber.Ctx, err error) error {
	var error *Error
	var fiberError *fiber.Error

	switch true {
	case errors.As(err, &error):
		break
	case errors.As(err, &fiberError):
		error.Status = fiberError.Code
		error.Title = fiberError.Message
	default:
		error.Status = fiber.StatusInternalServerError
		error.Code = "SERVER_ERROR"
		error.Detail = err.Error()
	}

	return ctx.Status(error.Status).JSON(error)
}
