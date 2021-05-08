package auth

import (
	"context"
	"errors"
)

var ErrAccessDenied = errors.New("access denied")

//maps to the authenticated user in context
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUserFromContext(ctx context.Context) *User {
	if ctx.Value(userCtxKey) == nil {
		return nil
	}
	return ctx.Value(userCtxKey).(*User)
}
