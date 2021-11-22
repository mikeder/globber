package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/mikeder/globber/internal/models"
)

const (
	accessTTL  = 15 * time.Minute
	refreshTTL = 24 * time.Hour
)

var signingMethod = jwt.SigningMethodHS256
var tokenCache = map[string]time.Time{}
var mu sync.Mutex

// Claims holds our authorized claims and standard JWT claims.
type Claims struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Email   string `json:"email"`
	jwt.RegisteredClaims
}

// Manager coordinates authentication and user methods.
type Manager struct {
	Auth          *jwtauth.JWTAuth
	signingSecret []byte
	name          string
	userDB        models.XODB
}

// NewManager returns a new instance of an authentication Manager.
func NewManager(secret []byte, name string, userDB models.XODB) *Manager {
	theMan := &Manager{
		Auth:          jwtauth.New(signingMethod.Alg(), secret, nil),
		signingSecret: secret,
		name:          name,
		userDB:        userDB,
	}

	return theMan
}

// ValidateCtx checks a context for a valid token and a username.
func ValidateCtx(ctx context.Context) (bool, string) {
	var user string

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		log.Print(errors.Wrap(err, "checking context for token and claims"))
		return false, user
	}

	if name, ok := claims["name"]; ok {
		user = name.(string)
	}

	return true, user
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
	claims := jwt.MapClaims{
		"aud": "globber-debug",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"sub": "superuser",
	}

	token := jwt.NewWithClaims(signingMethod, claims)
	tokenString, err := token.SignedString(m.signingSecret)
	if err != nil {
		log.Fatal(err)
	}

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

	// parse and validate incoming refresh token
	rt, err := jwt.Parse(t.Refresh, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return m.signingSecret, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "parse refresh token")
	}

	cl, ok := rt.Claims.(Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	mu.Lock()
	_, tokenKnown := tokenCache[cl.ID]
	mu.Unlock()

	if !tokenKnown {
		return nil, errors.New("unknown token")
	}

	uid, err := strconv.Atoi(cl.Subject)
	if err != nil {
		return nil, errors.Wrap(err, "parse subject uid")
	}

	user, err := models.AuthorByID(ctx, m.userDB, uid)
	if err != nil {
		return nil, errors.Wrap(err, "getting user from database")
	}

	return m.newTokens(user)
}

// ListTokens returns the current token cache list.
func (m *Manager) ListTokens(ctx context.Context) interface{} {
	type ret struct {
		Tokens map[string]time.Time `json:"tokens"`
	}

	mu.Lock()
	tokens := tokenCache
	mu.Unlock()

	return ret{tokens}
}

// Tokens contains access and refresh JWT's.
type Tokens struct {
	Access     string    `json:"access_token"`
	AccessTTL  time.Time `json:"access_ttl"`
	Refresh    string    `json:"refresh_token"`
	RefreshTTL time.Time `json:"refresh_ttl"`
}

func (m *Manager) newTokens(u *models.Author) (*Tokens, error) {
	now := time.Now()
	accessExp := now.Add(accessTTL)
	accessClaims := &Claims{
		Name:  u.Name,
		Email: u.Email,
	}

	refreshExp := now.Add(refreshTTL)
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"globber"},
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			ID:        uuid.New().String(),
			Issuer:    m.name,
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	var signErr error
	at := jwt.NewWithClaims(signingMethod, accessClaims)
	ats, signErr := at.SignedString(m.signingSecret)
	if signErr != nil {
		log.Fatal(signErr)
	}

	rt := jwt.NewWithClaims(signingMethod, refreshClaims)
	rts, signErr := rt.SignedString(m.signingSecret)
	if signErr != nil {
		log.Fatal(signErr)
	}

	mu.Lock()
	tokenCache[refreshClaims.ID] = time.Now()
	mu.Unlock()

	return &Tokens{
		Access:     ats,
		AccessTTL:  accessExp,
		Refresh:    rts,
		RefreshTTL: refreshExp,
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
