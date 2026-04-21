package queries

import (
	"context"
	"errors"
	"fmt"

	"github.com/marlon-clemente/timyo-playground-backend/packages/errs"
)

type WorkspaceReadModel struct {
	ID string `json:"id"`
	Name string `json:"name"`
	AgentID string `json:"agentID"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}


type MeReadModel struct {
	ID string `json:"id"`
	MemberID string `json:"memberID"`
	Name string `json:"name"`
	AvatarURL string `json:"avatarURL"`
	CreatedAt string `json:"createdAt"`
	Workspace *WorkspaceReadModel `json:"workspace,omitempty"`
	FormCount int `json:"formCount"`
}

type IMe interface {
	Me(ctx context.Context, memberID string) (*MeReadModel, error)
} 
	
type Me struct {
	query IMe
}

func NewMeQuery(query IMe) *Me {
	return &Me{query: query}
}

func (q *Me) Me(ctx context.Context, memberID string) (*MeReadModel, error) {
	data, err := q.query.Me(ctx, memberID)
	if err != nil {
		return nil, err
	}

	if data.MemberID == "" {
		return nil, errs.NotFoundErr("agents not found", errors.New("agents not found")).WithCode("REQUIRED_AGENT")
	}
	
	if data.Workspace == nil || data.Workspace.ID == "" {
		msg := fmt.Sprintf("workspace not found for member ID %s", memberID)
		return nil, errs.NotFoundErr(msg, errors.New(msg)).WithCode("REQUIRED_WORKSPACE")
	}

	return data, nil
}