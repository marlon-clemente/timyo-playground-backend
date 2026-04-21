package entities

import (
	"time"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
)

type Agent struct {
	ID vos.UUID
	MemberID vos.UUID
	Name vos.Name
	AvatarURL string
	CreatedAt time.Time
}

func NewAgent(memberId vos.UUID, name vos.Name, avatarURL string) *Agent {
	return &Agent{
		ID: vos.NewUUID(),
		MemberID: memberId,
		Name: name,
		AvatarURL: avatarURL,
		CreatedAt: time.Now(),
	}
}