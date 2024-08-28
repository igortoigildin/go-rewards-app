package context

import (
	"context"
	"errors"
	"net/http"

	user "github.com/igortoigildin/go-rewards-app/internal/entities/user"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *user.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetUser(r *http.Request) (*user.User, error) {
	user, ok := r.Context().Value(userContextKey).(*user.User)
	if !ok {
		return nil, errors.New("missing user value in request context")
	}
	return user, nil
}
