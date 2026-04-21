package entities

import (
	"time"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
)


type Form struct {
	Id 		vos.UUID
	WorkspaceID 	vos.UUID
	AgentID 	vos.UUID // Owner of the form, not necessarily the creator

	Name        vos.Name
	Description vos.Description
	
	isPublic    bool
	
	CreatedAt   time.Time
	UpdatedAt   time.Time
	
	CurrentVersion 	vos.UUID
	Versions        []FormVersion

}

func NewForm(
	workspaceID vos.UUID,
	agentID vos.UUID,
	name vos.Name,
	description vos.Description,
) *Form {
	now := time.Now()
	
	formID := vos.NewUUID()
	version := make([]FormVersion, 0)

	v := NewInitialFormVersion(formID)
	version = append(version, *v)
	form := &Form{
		Id:             formID,
		WorkspaceID:    workspaceID,
		AgentID:        agentID,
		CurrentVersion: v.Id,
		Name:           name,
		Description:    description,
		isPublic:       false,
		CreatedAt:      now,
		UpdatedAt:      now,
		Versions:       version,
	}
	return form
}