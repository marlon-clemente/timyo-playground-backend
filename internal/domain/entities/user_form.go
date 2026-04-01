package entities

import "time"

type UserForm struct {
	ID          string
	UserID      string
	Name        string
	Description string
	IsShared    bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
