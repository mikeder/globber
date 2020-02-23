package auth

import "fmt"

// ErrUserMissingField is used when a User is missing either the
// email or password field.
type ErrUserMissingField struct {
	Field string
}

func (e ErrUserMissingField) Error() string {
	return fmt.Sprintf("user missing %s", e.Field)
}
