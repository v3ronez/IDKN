package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/v3ronez/IDKN/internal/validator"
)

const (
	ScopeActivation    = "activation"
	ScopeAuthenticaton = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

type TokenModel struct {
	DB *sql.DB
}

func (t TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	return token, err
}

func ValidateTokenPlainText(v *validator.Validator, TokenPlainText string) {
	v.Check(TokenPlainText != "", "token", "must be provided")
	v.Check(len(TokenPlainText) == 26, "token", "must be 26 bytes long")
}

func (t TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES($1,$2,$3,$4)`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (t TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()
	_, err := t.DB.ExecContext(ctx, query, scope, userID)
	if err != nil {
		return err
	}
	return nil
}
