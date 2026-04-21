package adapters

import (
	"context"
	"database/sql"

	"github.com/marlon-clemente/timyo-playground-backend/internal/application/queries"
	"github.com/marlon-clemente/timyo-playground-backend/packages/database"
	"gorm.io/gorm"
)


type MeQuery struct {
	db *gorm.DB
}

func NewMeQuery() queries.IMe {
	return &MeQuery{db: database.Get()}
}

type memberRow struct {
	ID string `gorm:"column:id"`
	MemberID string `gorm:"column:member_id"`
	Name string `gorm:"column:name"`
	AvatarURL string `gorm:"column:avatar_url"`
	CreatedAt string `gorm:"column:created_at"`
	WorkspaceID sql.NullString `gorm:"column:workspace_id"`
	WorkspaceName sql.NullString `gorm:"column:workspace_name"`
	WorkspaceAgentID sql.NullString `gorm:"column:workspace_agent_id"`
	WorkspaceCreatedAt sql.NullString `gorm:"column:workspace_created_at"`
	WorkspaceUpdatedAt sql.NullString `gorm:"column:workspace_updated_at"`
	FormCount int `gorm:"column:form_count"`
}

func (q *MeQuery) Me(ctx context.Context, memberID string) (*queries.MeReadModel, error) {
	sqlQuery := `
		SELECT
			a.id,
			a.member_id,
			a.name,
			a.avatar_url,
			a.created_at,
			w.id AS workspace_id,
			w.name AS workspace_name,
			w.agent_id AS workspace_agent_id,
			w.created_at AS workspace_created_at,
			w.updated_at AS workspace_updated_at,
			COALESCE((
				SELECT COUNT(*)
				FROM forms f
				WHERE f.workspace_id = w.id
			), 0) AS form_count
		FROM agents a
		LEFT JOIN LATERAL (
			SELECT
				id,
				name,
				agent_id,
				created_at,
				updated_at
			FROM workspaces
			WHERE agent_id = a.id
			ORDER BY created_at DESC
			LIMIT 1
		) w ON TRUE
		WHERE a.member_id = ?
		LIMIT 1;
	`

	var me memberRow
	err := q.db.WithContext(ctx).Raw(sqlQuery, memberID).Scan(&me).Error
	if err != nil {
		return nil, err
	}

	var workspace *queries.WorkspaceReadModel
	if me.WorkspaceID.Valid {
		workspace = &queries.WorkspaceReadModel{
			ID: me.WorkspaceID.String,
			Name: me.WorkspaceName.String,
			AgentID: me.WorkspaceAgentID.String,
			CreatedAt: me.WorkspaceCreatedAt.String,
			UpdatedAt: me.WorkspaceUpdatedAt.String,
		}
	}

	return &queries.MeReadModel{
		ID:        me.ID,
		MemberID:  me.MemberID,
		Name:      me.Name,
		AvatarURL: me.AvatarURL,
		CreatedAt: me.CreatedAt,
		Workspace: workspace,
		FormCount: me.FormCount,
	}, nil
}
