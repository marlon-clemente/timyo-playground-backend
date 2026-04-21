package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/marlon-clemente/timyo-playground-backend/adapters"
	"github.com/marlon-clemente/timyo-playground-backend/interface/handlers"
	"github.com/marlon-clemente/timyo-playground-backend/internal/application/queries"
	usecases "github.com/marlon-clemente/timyo-playground-backend/internal/application/use-cases"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

func PrivateRoutes(app *fiber.App, jwtSecret string, authServiceURL string) {
	// Adapters
	userInfo := adapters.NewUserInfo(authServiceURL)

	agentRepo := adapters.NewAgentRepo()
	workspaceRepo := adapters.NewWorkspaceRepo()
	formRepo := adapters.NewFormDB()
	
	// queries
	meQuery := adapters.NewMeQuery()

	// use cases
	agentsUseCase := usecases.NewAgent(agentRepo, userInfo)
	workspaceUseCase := usecases.NewWorkspace(workspaceRepo)
	formUseCase := usecases.NewFormUseCase(formRepo)
	
	// querie handlers
	meQueryHandler := queries.NewMeQuery(meQuery)

	// handlers
	meHandler := handlers.NewMeHandler(*meQueryHandler)
	agentsHandler := handlers.NewAgentsHandler(*agentsUseCase)
	formHandler := handlers.NewFormsHandler(formUseCase)
	
	wpHandlers := handlers.NewWorkspaceHandler(workspaceUseCase)

	private := app.Group("", server.AuthMiddleware(jwtSecret))
	private.Use(server.ContextMiddleware(meQuery))



	private.Get("/me", server.Adapt(meHandler.GetMe))
	private.Post("/agent", server.Adapt(agentsHandler.CreateAgent))

	private.Post("/workspace", server.Adapt(wpHandlers.CreateWorkspace))

	private.Post("/forms", server.Adapt(formHandler.CreateForm))
	private.Get("/forms", server.Adapt(formHandler.ListForms))
	private.Get("/forms/:formId", server.Adapt(formHandler.GetForm))

	private.Put("/forms/:formId/version", server.Adapt(formHandler.SaveFormVersion))
	private.Delete("/forms/:formId", server.Adapt(formHandler.DeleteForm))
}