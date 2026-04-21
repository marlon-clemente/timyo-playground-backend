package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marlon-clemente/timyo-playground-backend/internal/application/queries"
)

// ContextMiddleware enriches request context with agent/workspace identifiers
// derived from the authenticated member ID.
func ContextMiddleware(meQuery queries.IMe) fiber.Handler {
	return func(c *fiber.Ctx) error {
		memberID, ok := c.Locals("userId").(string)
		if !ok || memberID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authenticated user",
			})
		}

		me, err := meQuery.Me(c.UserContext(), memberID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Unable to resolve user context",
				"details": err.Error(),
			})
		}

		if me != nil {
			if me.ID != "" {
				c.Locals("agentId", me.ID)
			}
			if me.Workspace != nil && me.Workspace.ID != "" {
				c.Locals("workspaceId", me.Workspace.ID)
			}
		}

		return c.Next()
	}
}