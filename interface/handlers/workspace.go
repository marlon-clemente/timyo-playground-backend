package handlers

import (
	usecases "github.com/marlon-clemente/timyo-playground-backend/internal/application/use-cases"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

type WorkspaceHandler struct {
	uc *usecases.Workspace
}

func NewWorkspaceHandler(uc *usecases.Workspace) *WorkspaceHandler {
	return &WorkspaceHandler{uc: uc}
}

func (h *WorkspaceHandler) CreateWorkspace(c *server.Ctx) error {
	var input WorkspaceCreateRequest
	if err := c.BindAndValidate(&input); err != nil {
		return err
	}

	agentID := c.GetAgentID()
	if agentID == "" {
		return c.Status(401).JSON(map[string]string{
			"error": "Agent context not found",
		})
	}

	result, err := h.uc.Create(c.UserContext(), usecases.WorkspaceCreateInput{
		AgentID: agentID,
		Name: input.Name,
	})
	if err != nil {
		return err
	}

	return c.ResponseOk(result)
}