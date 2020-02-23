package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikeder/globber/internal/models"
)

// tokens are good for 5 minutes
var tokenExpiry = time.Now().Add(5 * time.Minute)

// Claims holds our authorized claims and standard JWT claims.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Manager coordinates authentication and user methods.
type Manager struct {
	TokenAuth *jwtauth.JWTAuth
	userDB    *sql.DB
}

// NewManager returns a new instance of an authentication Manager.
func NewManager(userDB *sql.DB, secret string) *Manager {
	return &Manager{jwtauth.New("HS256", []byte(secret), nil), userDB}
}

// PasswordLogin performs password authentication of a user.
func (m *Manager) PasswordLogin(ctx context.Context, c *Credentials) (string, error) {
	if c.Email == "" {
		return "", ErrUserMissingField{"email"}
	}
	if c.Password == "" {
		return "", ErrUserMissingField{"password"}
	}

	user, err := models.AuthorByEmail(ctx, m.userDB, c.Email)
	if err != nil {
		return "", errors.Wrap(err, "getting user from database")
	}

	valid := checkPasswordHash(c.Password, user.HashedPassword)
	if !valid {
		return "", errors.New("password did not match")
	}

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: c.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: tokenExpiry.Unix(),
		},
	}

	_, tokenString, err := m.TokenAuth.Encode(claims)

	return tokenString, err
}

// AddUser adds a new User to the database.
func (m *Manager) AddUser(ctx context.Context, u *User) error {
	_, claims, err := jwtauth.FromContext(ctx)
	if caller, ok := claims["username"]; !ok || caller.(string) != "SuperUser" {
		return errors.New("not authorized to add user")
	}

	if u.Email == "" {
		return ErrUserMissingField{"email"}
	}
	if u.Name == "" {
		return ErrUserMissingField{"name"}
	}
	if u.Password == "" {
		return ErrUserMissingField{"password"}
	}

	author, err := models.AuthorByEmail(ctx, m.userDB, u.Email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if author != nil {
		return errors.New("email already in use")
	}

	hash, err := hashPassword(u.Password)
	if err != nil {
		return err
	}

	newUser := models.Author{
		Email:          u.Email,
		Name:           u.Name,
		HashedPassword: hash,
	}

	if err := newUser.Insert(ctx, m.userDB); err != nil {
		return err
	}

	return nil
}

// DebugToken returns a token for debug purposes, it is valid for 1 hour.
func (m *Manager) DebugToken() string {
	expires := time.Now().Add(time.Hour * 1)
	_, tokenString, _ := m.TokenAuth.Encode(Claims{
		Username: "SuperUser",
		StandardClaims: jwt.StandardClaims{
			Audience:  "debug",
			ExpiresAt: expires.Unix(),
		},
	})
	return tokenString
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
