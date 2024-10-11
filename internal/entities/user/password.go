package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Plaintext *string
	Hash      []byte
}

// Calculates the bcrypt hash of a plaintext password.
func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

// Checks whether the provided plaintext password matches the hashed password stores in the struct.
func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, err
}
