package api

import (
	"context"
	"errors"
	"net/http"

	entities "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type contextKey string

const userContextKey = contextKey("user")

func contextSetUser(r *http.Request, user *entities.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func contextGetUser(r *http.Request) (*entities.User, error) {
	user, ok := r.Context().Value(userContextKey).(*entities.User)
	if !ok {
		return nil, errors.New("missing user value in request context")
	}
	return user, nil
}
