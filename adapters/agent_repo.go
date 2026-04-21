package adapters

import (
	"context"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/entities"
	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/repositories"
	"github.com/marlon-clemente/timyo-playground-backend/packages/database"
	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
	"gorm.io/gorm"
)

type AgentRepo struct {
	db *gorm.DB
}

func NewAgentRepo() repositories.Agents {
	return &AgentRepo{db: database.Get()}
}

func (r *AgentRepo) FindByID(ctx context.Context, id string) (*entities.Agent, error) {
	sql := `SELECT id, member_id, name, avatar_url, created_at FROM agents WHERE id = ? LIMIT 1;`
	var agent entities.Agent
	tx := r.db.WithContext(ctx).Raw(sql, id).Scan(&agent)
	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected == 0 {
		return nil, errs.NotFoundErr("agent not found", gorm.ErrRecordNotFound).WithCode("REQUIRED_AGENT")
	}

	return &agent, nil	
}
func (r *AgentRepo) Create(ctx context.Context, agent *entities.Agent) error {
	sql := `INSERT INTO agents (id, member_id, name, avatar_url, created_at) VALUES (?, ?, ?, ?, ?);`
	return r.db.WithContext(ctx).Exec(sql, agent.ID, agent.MemberID, agent.Name, agent.AvatarURL, agent.CreatedAt).Error
}

func (r *AgentRepo) Update(ctx context.Context, agent *entities.Agent) error {
	sql := `UPDATE agents SET member_id = ?, name = ?, avatar_url = ? WHERE id = ?;`
	return r.db.WithContext(ctx).Exec(sql, agent.MemberID, agent.Name, agent.AvatarURL, agent.ID).Error
}