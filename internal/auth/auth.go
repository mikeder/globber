package auth

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikeder/globber/internal/models"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 24 * time.Hour
	tokenIss   = "sqweeb.net"
)

var tokenCache = map[string]struct{}{}
var mu sync.Mutex

// Claims holds our authorized claims and standard JWT claims.
type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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

// ValidateCtx checks a context for a valid token and a username.
func ValidateCtx(ctx context.Context) (bool, string) {
	var user string

	token, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		log.Print(errors.Wrap(err, "checking context for token and claims"))
		return false, user
	}

	if name, ok := claims["name"]; ok {
		user = name.(string)
	}

	return token.Valid, user
}

// AddUser adds a new User to the database.
func (m *Manager) AddUser(ctx context.Context, u *User) error {
	// validate incoming user fields
	if err := u.validate(); err != nil {
		return err
	}

	// validate incoming auth token
	if valid, username := ValidateCtx(ctx); !valid || username != "superuser" {
		return errors.New("not authorized to add user")
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
		Name: "superuser",
		StandardClaims: jwt.StandardClaims{
			Audience:  "debug",
			ExpiresAt: expires.Unix(),
		},
	})
	return tokenString
}

// PasswordLogin performs password authentication of a user.
func (m *Manager) PasswordLogin(ctx context.Context, c *Credentials) (tokens *Tokens, err error) {
	if c.Email == "" {
		return nil, ErrUserMissingField{"email"}
	}
	if c.Password == "" {
		return nil, ErrUserMissingField{"password"}
	}

	user, err := models.AuthorByEmail(ctx, m.userDB, c.Email)
	if err != nil {
		return nil, errors.Wrap(err, "getting user from database")
	}

	valid := checkPasswordHash(c.Password, user.HashedPassword)
	if !valid {
		return nil, errors.New("password did not match")
	}

	return m.newTokens(user)
}

// Refresh performs a token check and issues new tokens if valid.
func (m *Manager) Refresh(ctx context.Context, t *Tokens) (*Tokens, error) {
	if t == nil {
		return nil, errors.New("nil token given to refresh")
	}

	validToken, err := m.TokenAuth.Decode(t.Refresh.Raw)
	if err != nil {
		return nil, err
	}

	var claims jwt.MapClaims
	if tokenClaims, ok := validToken.Claims.(jwt.MapClaims); ok {
		claims = tokenClaims
	}

	suid := claims["sub"].(string)

	uid, err := strconv.ParseInt(suid, 10, 64)
	if err != nil {
		return nil, errors.New("bad subject")
	}

	log.Print("perform further token validation")
	log.Print(validToken.Valid)

	user, err := models.AuthorByID(ctx, m.userDB, int(uid))
	if err != nil {
		return nil, errors.Wrap(err, "getting user from database")
	}

	return m.newTokens(user)
}

// Tokens contains access and refresh JWT's.
type Tokens struct {
	Access  *jwt.Token `json:"access_token"`
	Refresh *jwt.Token `json:"refresh_token"`
}

func (m *Manager) newTokens(u *models.Author) (*Tokens, error) {
	now := time.Now()
	accessClaims := &Claims{
		Name:  u.Name,
		Email: u.Email,
		StandardClaims: jwt.StandardClaims{
			Audience:  "globber",
			ExpiresAt: now.Add(accessTTL).Unix(),
			Id:        uuid.New().String(),
			Issuer:    tokenIss,
			IssuedAt:  now.Unix(),
		},
	}

	refreshClaims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(accessTTL).Unix(),
			Subject:   string(u.ID),
			Id:        uuid.New().String(),
			Issuer:    tokenIss,
			IssuedAt:  now.Unix(),
		},
	}

	access, _, err := m.TokenAuth.Encode(accessClaims)
	if err != nil {
		return nil, err
	}
	refresh, _, err := m.TokenAuth.Encode(refreshClaims)
	if err != nil {
		return nil, err
	}

	mu.Lock()
	tokenCache[refreshClaims.Id] = struct{}{}
	mu.Unlock()

	return &Tokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
