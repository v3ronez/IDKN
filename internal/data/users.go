package data

import (
	"errors"
	"time"

	"github.com/v3ronez/IDKN/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}
	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func validatePaswordPlainText(v *validator.Validator, plainTextPassword string) {
	v.Check(plainTextPassword != "", "password", "must be provided")
	v.Check(len(plainTextPassword) >= 8, "password", "must at least 8 bytes long")
	v.Check(len(plainTextPassword) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "mest be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	ValidateEmail(v, user.Email)

	if user.Password.plainText != nil {
		validatePaswordPlainText(v, *user.Password.plainText)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
