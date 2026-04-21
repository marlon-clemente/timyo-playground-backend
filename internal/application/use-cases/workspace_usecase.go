package usecases

import (
	"context"
	"fmt"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/repositories"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
)

type Workspace struct {
	repo repositories.Workspaces	
}

func NewWorkspace(repo repositories.Workspaces) *Workspace {
	return &Workspace{repo: repo}
}

func (w *Workspace) Create(ctx context.Context, input WorkspaceCreateInput) (WorkspaceOutput, error) {
	name, err := vos.NewName(input.Name)
	if err != nil {
		return WorkspaceOutput{}, err
	}

	id, err := vos.ParseUUID(input.AgentID)
	if err != nil {
		return WorkspaceOutput{}, err
	}

	wk := entities.NewWorkspace(id, name)

	err = w.repo.Create(ctx, wk)
	if err != nil {
		return WorkspaceOutput{}, fmt.Errorf("failed to create workspace: %w", err)
	}

	return WorkspaceOutput{
		ID:        wk.ID.String(),
		Name:      wk.Name.String(),
		AgentID:   wk.AgentID.String(),
		CreatedAt: wk.CreatedAt.String(),
		UpdatedAt: wk.UpdatedAt.String(),
	}, nil
}