package server

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Ctx is a wrapper around *fiber.Ctx to provide common helper methods
type Ctx struct {
	*fiber.Ctx
	Logger *slog.Logger
}

// Log returns the contextual logger. If no logger is set, it returns the default slog logger.
func (c *Ctx) Log() *slog.Logger {
	if c.Logger == nil {
		return slog.Default()
	}
	return c.Logger
}

// BindAndValidate parses the request body into the provided struct and validates it using struct tags.
// It returns an error if either parsing or validation fails.
func (c *Ctx) BindAndValidate(out any) error {
	if err := c.BodyParser(out); err != nil {
		return errs.DomainErr("invalid request body", err)
	}
	if err := validate.Struct(out); err != nil {
		var fields []map[string]any
		for _, fe := range err.(validator.ValidationErrors) {
			fields = append(fields, map[string]any{
				"field":  fe.Field(),
				"reason": fe.Tag(),
				"type":   fe.Type().String(),
			})
		}
		return &errs.Error{
			Message: "validation failed",
			Err:     err,
			Variant: errs.ErrValidation,
			Data:    map[string]any{"invalidFields": fields},
		}
	}
	return nil
}

// GetUserID retrieves the user authentication ID from context
func (c *Ctx) GetUserID() string {
	userID, ok := c.Locals("userId").(string)
	if !ok {
		return ""
	}
	return userID
}

func (c *Ctx) GetAgentID() string {
	agentID, ok := c.Locals("agentId").(string)
	if !ok {
		return ""
	}
	return agentID
}

func (c *Ctx) GetWorkspaceID() string {
	workspaceID, ok := c.Locals("workspaceId").(string)
	if !ok {
		return ""
	}
	return workspaceID
}

// NewLoggingMiddleware returns a Fiber middleware that injects the provided slog.Logger into the request context.
func NewLoggingMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("logger", log)
		return c.Next()
	}
}

// Handler is a function type that takes our custom Ctx
type Handler func(c *Ctx) error

// Adapt converts a server.Handler to a fiber.Handler.
// It retrieves the logger from the request context ("logger" local) and passes it to the custom Ctx.
func Adapt(h Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log, _ := c.Locals("logger").(*slog.Logger)
		return h(&Ctx{Ctx: c, Logger: log})
	}
}
