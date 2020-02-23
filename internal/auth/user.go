package auth

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
