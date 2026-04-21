package handlers

import (
	usecases "github.com/marlon-clemente/timyo-playground-backend/internal/application/use-cases"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

type AgentsHandler struct {
	agentsUseCase usecases.Agent
}

func NewAgentsHandler(agentsUseCase usecases.Agent) *AgentsHandler {
	return &AgentsHandler{agentsUseCase: agentsUseCase}
}

func (h *AgentsHandler) CreateAgent(c *server.Ctx) error {
	
	accessToken := c.Locals("accessToken")

	input := usecases.CreateAgentInput{
		MemberID: c.GetUserID(),
		AccessToken: accessToken.(string),
	}

	output, err := h.agentsUseCase.Create(c.UserContext(), input)
	if err != nil {
		return err
	}
	return c.ResponseOk(output)
}