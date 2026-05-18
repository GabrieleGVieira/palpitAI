package apperrors

import "errors"

var (
	ErrConflict   = errors.New("conflict")
	ErrForbidden  = errors.New("forbidden")
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation")
)

type Error struct {
	Kind    error
	Message string
}

func (err Error) Error() string {
	return err.Message
}

func (err Error) Is(target error) bool {
	return target == err.Kind
}

func NewConflict(message string) error {
	return Error{Kind: ErrConflict, Message: message}
}

func NewForbidden(message string) error {
	return Error{Kind: ErrForbidden, Message: message}
}

func NewNotFound(message string) error {
	return Error{Kind: ErrNotFound, Message: message}
}

func NewValidation(message string) error {
	return Error{Kind: ErrValidation, Message: message}
}

func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}
