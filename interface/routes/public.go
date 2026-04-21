package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/marlon-clemente/timyo-playground-backend/interface/handlers"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

const scalarHTML = `
<!DOCTYPE html>
<html>
  <head>
    <title>Timyo Playground API - Scalar</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <!-- Need a theme? https://sandbox.scalar.com/theme -->
    <script
      id="api-reference"
      data-url="/swagger/doc.json"
    ></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`


func PublicRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Scalar UI serving the swagger definition
	app.Get("/scalar", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(scalarHTML)
	})

	app.Get("/health", server.Adapt(handlers.HealthCheck))
	app.Get("/ping", server.Adapt(handlers.Ping))


}
