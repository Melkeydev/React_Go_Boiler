package main


// This is going to hold our context for our User (if anonymous or activated)
import (
	"context"
	"net/http"
	"backend/models"
)

type contextKey string

// store the value of the token
const userContextKey = contextKey("user")

// setUserContext
func (app *application) contextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}  

// TODO: need to investigate this
//getUserContext
func (app *application) contextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)

	if !ok {
		panic("missing user value in request context")
	}

	return user
}
