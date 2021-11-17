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
	

	user, err := app.models.DB.
}
