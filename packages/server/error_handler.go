package server

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
	"github.com/marlon-clemente/timyo-playground-backend/packages/observability"
)

var errVariantToStatus = map[errs.ErrsType]int{
	errs.ErrUnauthorized: http.StatusUnauthorized,
	errs.ErrForbidden:    http.StatusForbidden,
	errs.ErrValidation:   http.StatusUnprocessableEntity,
	errs.ErrNotFound:     http.StatusNotFound,
	errs.ErrDomain:       http.StatusBadRequest,
	errs.ErrConflict:     http.StatusConflict,
	errs.ErrInternal:     http.StatusInternalServerError,
}

type ErrsResponse struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
	Code    string         `json:"code,omitempty"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	ctx := c.UserContext()

	// Domain/application errors from the errs package
	if appErr, ok := errors.AsType[*errs.Error](err); ok {
		status, ok := errVariantToStatus[appErr.Variant]
		if !ok {
			status = http.StatusInternalServerError
		}

		// Send the real underlying cause to the tracer, not the user-facing message
		realErr := appErr.Err
		if realErr == nil {
			realErr = appErr
		}
		observability.AddError(ctx, realErr)

		return c.Status(status).JSON(ErrsResponse{
			Message: appErr.Message,
			Data:    appErr.Data,
			Code:    appErr.Code,
		})
	}

	// Fiber's own errors (e.g. 404 from unmatched route)
	if fiberErr, ok := errors.AsType[*fiber.Error](err); ok {
		observability.AddError(ctx, fiberErr)
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"message": fiberErr.Message,
		})
	}

	// Fallback: unexpected internal error
	observability.AddError(ctx, err)
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"message": "internal server error",
	})
}
