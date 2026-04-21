package entities

import (
	"time"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
)

type Workspace struct {
	ID   vos.UUID
	AgentID vos.UUID
	Name vos.Name
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWorkspace(agentId vos.UUID, name vos.Name) *Workspace {
	now := time.Now()
	return &Workspace{
		ID: vos.NewUUID(),
		AgentID: agentId,
		Name: name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *Workspace) UpdateName(name vos.Name) {
	w.Name = name
	w.UpdatedAt = time.Now()
}