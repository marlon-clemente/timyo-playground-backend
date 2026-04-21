package entities

import (
	"time"

	"github.com/marlon-clemente/timyo-playground-backend/internal/domain/vos"
)

type FormVersion struct {
	Id 		vos.UUID
	FormID 	vos.UUID
	
	VersionNumber int

	Props any
	
	CreatedAt   time.Time
	UpdatedAt   time.Time

}

func NewInitialFormVersion(formID vos.UUID) *FormVersion {
	now := time.Now()
	return &FormVersion{
		Id:            vos.NewUUID(),
		FormID:        formID,
		VersionNumber: 1,
		Props:         nil, // Initial version has no props
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}