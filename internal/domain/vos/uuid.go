package vos

import (
	"database/sql/driver"
	"errors"

	"github.com/google/uuid"
)

type UUID string

func NewUUID() UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return UUID(id.String())
}

func ParseUUID(value string) (UUID, error) {
	if _, err := uuid.Parse(value); err != nil {
		return "", errors.New("invalid UUID")
	}
	return UUID(value), nil
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) Value() (driver.Value, error) {
	if u == "" {
		return nil, nil
	}
	return string(u), nil
}

func (u *UUID) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*u = UUID(v)
		return nil
	case []byte:
		*u = UUID(v)
		return nil
	case nil:
		*u = ""
		return nil
	}
	return errors.New("unsupported type for UUID")
}
