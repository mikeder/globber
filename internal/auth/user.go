package auth

import (
	"regexp"

	"github.com/pkg/errors"
)

// Credentials represents a user performing an authentication request.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a new user.
type User struct {
	Name string `json:"name"`
	Credentials
}

func (u *User) validate() error {
	if u.Email == "" {
		return ErrUserMissingField{"email"}
	}
	if u.Name == "" {
		return ErrUserMissingField{"name"}
	}
	if u.Password == "" {
		return ErrUserMissingField{"password"}
	}

	if len(u.Password) < 8 || len(u.Password) > 128 {
		return errors.New("password must be between 8 and 128 characters long")
	}

	if !validEmail(u.Email) {
		return errors.New("invalid email address")
	}
	return nil
}

func validEmail(address string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(address) > 254 || !rxEmail.MatchString(address) {
		return false
	}
	return true
}
