package main

import (
	"backend/models"
	"backend/validator"
	"errors"
	"net/http"
	"time"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the email and password
	v := validator.New()
	
	models.ValidateEmail(v, input.Email)
	models.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.DB.GetUserByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.invalidCredentialResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	
	// Compare and match the hash passwords
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// if passwords do not match
	if !match {
		app.invalidCredentialResponse(w, r)
		return
	}

	// Creates a new auth token and saves it
	token, err := app.models.DB.NewToken(user.ID, 24*time.Hour, models.ScopeAuthentication) 
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token":token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}















