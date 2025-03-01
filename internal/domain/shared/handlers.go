package shared

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandler() func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		var response Error

		var fiberError *fiber.Error

		switch true {
		case errors.As(err, &fiberError):
			response.Status = fiberError.Code
			// switch code := fiberError.Code; {
			// case code > 400:
			//
			// }
			// response.Code = strings.ReplaceAll(strings.ToUpper(fiberError.Message), " ", "_")
			response.Title = fiberError.Message
		default:
			response.Status = fiber.StatusInternalServerError
			response.Code = "INTERNAL_SERVER_ERROR"
			response.Detail = err.Error()
		}

		return ctx.Status(response.Status).JSON(response)
	}
}
