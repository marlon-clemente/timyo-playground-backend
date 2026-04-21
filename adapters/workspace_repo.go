package adapters

import (
	"context"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/repositories"
	"github.com/marlon-clemente/timyo-playground-backend/packages/database"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
	"gorm.io/gorm"
)

type workspaceRepo struct {
	db *gorm.DB
}

func NewWorkspaceRepo() repositories.Workspaces {
	return &workspaceRepo{db: database.Get()}
}

func (r *workspaceRepo) FindByID(ctx context.Context, id string) (*entities.Workspace, error) {
	sql := `SELECT id, name, agent_id, created_at, updated_at FROM workspaces WHERE id = ? LIMIT 1;`
	var workspace entities.Workspace
	tx := r.db.WithContext(ctx).Raw(sql, id).Scan(&workspace)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, errs.NotFoundErr("workspace not found", gorm.ErrRecordNotFound).WithCode("REQUIRED_WORKSPACE")
	}

	return &workspace, nil	
}

func (r *workspaceRepo) Create(ctx context.Context, workspace *entities.Workspace) error {
	sql := `INSERT INTO workspaces (id, name, agent_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?);`
	return r.db.WithContext(ctx).Exec(sql, workspace.ID, workspace.Name, workspace.AgentID, workspace.CreatedAt, workspace.UpdatedAt).Error
}

func (r *workspaceRepo) Update(ctx context.Context, workspace *entities.Workspace) error {
	sql := `UPDATE workspaces SET name = ?, agent_id = ?, updated_at = ? WHERE id = ?;`
	return r.db.WithContext(ctx).Exec(sql, workspace.Name, workspace.AgentID, workspace.UpdatedAt, workspace.ID).Error
}

