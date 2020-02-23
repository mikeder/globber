package auth

// Role represents a user's access level
type Role string

const (
	RoleAdmin  = Role("admin")
	RoleAuthor = Role("author")
	RoleUser   = Role("user")
	RoleNobody = Role("nobody")
)
