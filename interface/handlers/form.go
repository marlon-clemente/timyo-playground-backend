package handlers

import (
	"fmt"

	usecases "github.com/marlon-clemente/timyo-playground-backend/internal/application/use-cases"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"
)

type FormsHandler struct {
	uc *usecases.Form
}

func NewFormsHandler(uc *usecases.Form) *FormsHandler {
	return &FormsHandler{uc: uc}
}

func (h *FormsHandler) CreateForm(c *server.Ctx) error {
	var input FormCreateRequest
	if err := c.BindAndValidate(&input); err != nil {
		return err
	}

	agentID := c.GetAgentID()
	if agentID == "" {
		return c.Status(401).JSON(map[string]string{
			"error": "Agent context not found",
		})
	}

	workspaceID := c.GetWorkspaceID()
	if workspaceID == "" {
		return c.Status(401).JSON(map[string]string{
			"error": "Workspace context not found",
		})
	}

	result, err := h.uc.Create(c.UserContext(), usecases.CreateFormInput{
		AgentID:     agentID,
		Name:        input.Name,
		Description: input.Description,
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return err
	}

	return c.ResponseOk(result)
}

func (h *FormsHandler) ListForms(c *server.Ctx) error {
	workspaceID := c.GetWorkspaceID()
	if workspaceID == "" {
		return c.Status(401).JSON(map[string]string{
			"error": "Workspace context not found",
		})
	}

	result, err := h.uc.ListAlls(c.UserContext(), usecases.ListFormsInput{
		WorkspaceID: workspaceID,
	})
	if err != nil {
		return err
	}

	return c.ResponseOk(result)
}

func (h *FormsHandler) GetForm(c *server.Ctx) error {
	formID := c.Params("formId")


	if formID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "Form ID is required",
		})
	}

	result, err := h.uc.GetByID(c.UserContext(), formID)
	if err != nil {
		return err
	}

	return c.ResponseOk(result)
}

func (h *FormsHandler) SaveFormVersion(c *server.Ctx) error {
	formID := c.Params("formId")
	if formID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "Form ID is required",
		})
	}

	var input SaveFormVersionRequest
	if err := c.BindAndValidate(&input); err != nil {
		return err
	}

	_, err := h.uc.SaveFormVersion(c.UserContext(), usecases.SaveFormVersionInput{
		FormID: formID,
		Props: input.Props,
	})

	if err != nil {
		return fmt.Errorf("failed to save form version: %w", err)
	}

	return c.ResponseNoContent()
}

func (h *FormsHandler) DeleteForm(c *server.Ctx) error {
	formID := c.Params("formId")
	if formID == "" {
		return c.Status(400).JSON(map[string]string{
			"error": "Form ID is required",
		})
	}

	workspaceID := c.GetWorkspaceID()
	if workspaceID == "" {
		return c.Status(401).JSON(map[string]string{
			"error": "Workspace context not found",
		})
	}

	err := h.uc.Delete(c.UserContext(), usecases.DeleteFormInput{
		WorkspaceID: workspaceID,
		FormID: formID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete form: %w", err)
	}

	return c.ResponseNoContent()
}