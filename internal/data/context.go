package data

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func ContextSetUser(r *http.Request, user *User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetUser(ctx context.Context) *User {
	fmt.Println(ctx)
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
