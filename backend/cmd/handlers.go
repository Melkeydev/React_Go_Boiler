package main

import (
	"backend/models"
	"backend/validator"
	"encoding/json"
	"errors"
	//"fmt"
	//"github.com/pascaldekloe/jwt"
	//"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	//"time"
)

// Create a JSON message struct
type JSONMessage struct {
	Message string `json:"message"`
}

// Create a register user type
type UserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create a generic DBLoad type
type DBLoadPayload struct {
	DBDataOne   string `json:"db_data_one"`
	DBDataTwo   string `json:"db_data_two"`
	DBDataThree string `json:"db_data_three"`
}

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Status string
	}{"Curling Request"}

	js, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// generate the password hash with bcrypt
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if models.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DB.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user":user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// This function is client side - to - database
func (app *application) getData(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	data, err := app.models.DB.GetData(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": data}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// TODO: Refactor this aswell
//func (app *application) login(w http.ResponseWriter, r *http.Request) {
	//var payload UserPayload

	//err := json.NewDecoder(r.Body).Decode(&payload)
	//if err != nil {
		//app.logger.PrintError(err, nil)
	//}

	////we need to get the user
	//// TODO: replace with get usr by email
	//user, err := app.models.DB.GetUser(payload.Username)
	//if err != nil {
		//app.logger.PrintInfo("User does not exist", nil)
		//return
	//}

	//hashPassword := user.Password

	//err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(payload.Password))
	//// Handle the error for hasing and comparing
	//if err != nil {
		//log.Println(err)
		//_message := JSONMessage{
			//Message: "Unauthorized",
		//}

		//js, err := json.MarshalIndent(_message, "", "\t")
		//if err != nil {
			//app.logger.PrintError(err, nil)
		//}

		//w.Header().Set("Context-Type", "application/json")
		//w.WriteHeader(http.StatusOK)
		//w.Write(js)
		//return
	//}

	//// Validating a users token
	//var claims jwt.Claims
	//claims.Subject = fmt.Sprint(user.ID)
	//claims.Issued = jwt.NewNumericTime(time.Now())
	//claims.NotBefore = jwt.NewNumericTime(time.Now())
	//claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	//// supposed to be a unique domain you own
	//claims.Issuer = "github.com/melkeydev"
	//claims.Audiences = []string{"github.com/melkeydev"}

	//jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.Jwt.Secret))
	//if err != nil {
		//fmt.Println(err)
		//message := "Could not generate proper access"
		//app.errorResponse(w, r, http.StatusInternalServerError, message)
		//return
	//}

	////app.writeJSON(w, http.StatusOK, string(jwtBytes), "Successfully logged in")
	//_message := JSONMessage{
		//Message: string(jwtBytes),
	//}

	//js, err := json.MarshalIndent(_message, "", "\t")
	//if err != nil {
		//app.logger.PrintError(err, nil)
	//}

	//w.Header().Set("Context-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//w.Write(js)

//}

func (app *application) insertPayload(w http.ResponseWriter, r *http.Request) {
	var payload DBLoadPayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		return
	}

	dbload := &models.DBLoad{
		DBDataOne:   payload.DBDataOne,
		DBDataTwo:   payload.DBDataTwo,
		DBDataThree: payload.DBDataThree,
	}

	v := validator.New()

	if models.ValidateDBLoad(v, dbload); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DB.InsertDBLoad(dbload)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": dbload}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteDBload(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.DB.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "data deleted Succesfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateDBData(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// This is where we pull all the data first to update
	data, err := app.models.DB.GetData(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		DBDataOne   *string `json:"db_data_one"`
		DBDataTwo   *string `json:"db_data_two"`
		DBDataThree *string `json:"db_data_three"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Explicitly check each input
	if input.DBDataOne != nil {
		data.DBDataOne = *input.DBDataOne
	}

	if input.DBDataTwo != nil {
		data.DBDataTwo = *input.DBDataTwo
	}

	if input.DBDataThree != nil {
		data.DBDataThree = *input.DBDataThree
	}

	data.ID = id

	v := validator.New()

	// validate the json data
	if models.ValidateDBLoad(v, data); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DB.Update(data)
	if err != nil {
		switch {
		// the race condition editing error message
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": data}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listAllDBData(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DBDataOne string
		models.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	// Query string readers go below
	input.DBDataOne = app.readString(qs, "dbdataone", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "DBDataOne"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// We need to get all

	DBdata, metadata, err := app.models.DB.GetAll(input.DBDataOne, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"DBdata": DBdata, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
