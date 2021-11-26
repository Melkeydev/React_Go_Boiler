package main

import (
	"backend/models"
	"backend/validator"
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// We did not fully test this
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

// route and handle panics better
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})

}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Second {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.Limiter.Enabled {
			// get the ip addy from each request
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.Limiter.Rps), app.config.Limiter.Burst),
				}
			}

			// Every new ip that gets added to our clients slice gets a time stamp
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// Create a authenticate middleware
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adding ther VARY and AUTHORIZATION
		// This indicates to any cache that the request may vary
		w.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header
		authorizationHeader := r.Header.Get("Authorization")

		// The pointer reference might be questionable
		// IF there is no authorizationHeader we will set the context
		// this will hold an AnonymousUser - gifting bare minimum
		if authorizationHeader == "" {
			r = app.contextSetUser(r, models.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialResponse(w, r)
			return
		}

		token := headerParts[1]
		v := validator.New()

		// We need to validate the token to make sure it is correct format
		if models.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidCredentialResponse(w, r)
			return
		}

		// then we need to get the user
		user, err := app.models.DB.GetForToken(models.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, models.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
      return
		}

    // Set the user here pointer ref is questionable
    r = app.contextSetUser(r, user)

    // Because this is a MW wrapperm we need to pass to the next http handler
    next.ServeHTTP(w, r)
	})
}

// We need to split our auth to handle activated routes and authenticated routes
func(app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		// If this is true; You arent authorized
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
} 

// We need this to wrap and call our requireAuthenticatedUser MW
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w,r )
	})

	return app.requireAuthenticatedUser(fn)
}
