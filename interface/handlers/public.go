package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Check if the API is running
// @Tags         Public
// @Accept       json
// @Produce      json
// @Success      200  {object} server.Response
// @Failure      500  {object} server.Response
// @Router       /health [get]
func HealthCheck(c *server.Ctx) error {
	return c.ResponseOk(fiber.Map{
		"status": "ok",
	})
}

// Ping godoc
// @Summary      Ping
// @Description  Check if the API is running
// @Tags         Public
// @Accept       json
// @Produce      json
// @Success      200  {object} server.Response
// @Failure      500  {object} server.Response
// @Router       /ping [get]
func Ping(c *server.Ctx) error {
	return c.ResponseOk(fiber.Map{
		"status": "ok",
	})
}
