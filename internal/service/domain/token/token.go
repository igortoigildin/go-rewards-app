package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	userEntity "github.com/igortoigildin/go-rewards-app/internal/entities/user"
	"github.com/igortoigildin/go-rewards-app/internal/storage"
)

type TokenService struct {
	TokenRepository storage.TokenRepository
}

func (t *TokenService) NewToken(ctx context.Context, userID int64, ttl time.Duration) (*userEntity.Token, error) {
	token, err := generateToken(userID, ttl)
	if err != nil {
		return nil, err
	}

	err = t.TokenRepository.Insert(ctx, token)
	return token, err
}

func generateToken(useID int64, ttl time.Duration) (*userEntity.Token, error) {
	token := &userEntity.Token{
		UserID: useID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	// Fill the byte slice with random bytes from the operating system's CSPRNG
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}

// NewTokenService returns a new instance of user service.
func NewTokenService(TokenRepository storage.TokenRepository) *TokenService {
	return &TokenService{
		TokenRepository: TokenRepository,
	}
}
