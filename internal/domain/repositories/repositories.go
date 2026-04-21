package repositories

import (
	"context"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
)

type Agents interface {
	FindByID(ctx context.Context, id string) (*entities.Agent, error)
	Create(ctx context.Context, agent *entities.Agent) error
	Update(ctx context.Context, agent *entities.Agent) error
}

type Workspaces interface {
	FindByID(ctx context.Context, id string) (*entities.Workspace, error)
	Create(ctx context.Context, workspace *entities.Workspace) error
	Update(ctx context.Context, workspace *entities.Workspace) error
}

type Forms interface {
	Count(ctx context.Context, workspaceID string) (int, error)
	Create(ctx context.Context, form *entities.Form)  error
	FindByID(ctx context.Context, id string) (*entities.Form, error)
	ListByWorkspaceID(ctx context.Context, workspaceID vos.UUID) ([]*entities.Form, error)
	Delete(ctx context.Context, workspaceID, formID vos.UUID) error


	GetByIDWithVersion(ctx context.Context, id vos.UUID) (*entities.Form, error)
	UpdateFormVersion(ctx context.Context, version *entities.FormVersion) error
}