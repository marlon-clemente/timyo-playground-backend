package vos

import (
	"database/sql/driver"
	"errors"
	"regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email string

func NewEmail(value string) (Email, error) {
	if !emailRegex.MatchString(value) {
		return "", errors.New("invalid email address")
	}
	return Email(value), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Value() (driver.Value, error) {
	if e == "" {
		return nil, nil
	}
	return string(e), nil
}

func (e *Email) Scan(src any) error {
	switch v := src.(type) {
	case string:
		*e = Email(v)
		return nil
	case []byte:
		*e = Email(v)
		return nil
	case nil:
		*e = ""
		return nil
	}
	return errors.New("unsupported type for Email")
}