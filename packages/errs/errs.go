package errs

import "strings"

type ErrsType string

const (
	ErrUnauthorized ErrsType = "UNAUTHORIZED_ERROR"
	ErrForbidden    ErrsType = "FORBIDDEN_ERROR"
	ErrValidation   ErrsType = "VALIDATION_ERROR"
	ErrNotFound     ErrsType = "NOT_FOUND_ERROR"
	ErrDomain       ErrsType = "DOMAIN_ERROR"
	ErrConflict     ErrsType = "CONFLICT_ERROR"
	ErrInternal     ErrsType = "INTERNAL_ERROR"
)

type Error struct {
	// Mensagem descritiva para o usuário
	Message string         `json:"message"`
	Err     error          `json:"-"`
	Variant ErrsType       `json:"-"`
	Data    map[string]any `json:"data,omitempty"`
	Code   string         `json:"code,omitempty"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func UnauthorizedErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrUnauthorized,
	}
}

func NotFoundErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrNotFound,
	}
}

func ForbiddenErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrForbidden,
	}
}

func ConflictErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrConflict,
	}
}

func DomainErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrDomain,
	}
}

func ValidationErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrValidation,
	}
}

func InternalErr(message string, err error) *Error {
	return &Error{
		Message: message,
		Err:     err,
		Variant: ErrInternal,
	}
}

func (e *Error) AddAttribute(key string, value any) *Error {
	if e.Data == nil {
		e.Data = make(map[string]any)
	}
	e.Data[key] = value
	return e
}

func (e *Error) WithCode(code string) *Error {
	e.Code = strings.ToUpper(code)
	return e
}