package main

import (
  "fmt"
  "log"
  "time"
  "errors"
  "net/http"
  "encoding/json"
  "backend/models"
  "golang.org/x/crypto/bcrypt"
  "backend/validator"
  "github.com/pascaldekloe/jwt"
  //"github.com/julienschmidt/httprouter"
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
  DBDataOne string `json:"db_data_one"`
  DBDataTwo string `json:"db_data_two"`
  DBDataThree string `json:"db_data_three"`
}

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
  response := struct {
    Status string
  }{"Curling Request"}

  js, err := json.MarshalIndent(response, "", "\t")
  if err != nil {
    app.logger.Println(err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(js)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
  var payload UserPayload

  err := json.NewDecoder(r.Body).Decode(&payload)
  if err != nil {
    log.Println(err)
    return
  }

  var user models.User

  //hash paswords right away
  hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
  if err != nil {
    app.logger.Println(err)
  }

  user.Username = payload.Username
  user.Password = string(hashPassword)

  //Actual intereact with the DB
  err = app.models.DB.RegisterUser(user)
  if err != nil {
    app.logger.Println(err)
  }

  // Uses JSON message struct
  _message := JSONMessage{
    Message: "Succesfully registered a user", 
  } 

  js, err := json.MarshalIndent(_message, "", "\t")
  if err != nil {
    app.logger.Println(err)
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(js)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
  var payload UserPayload

  err := json.NewDecoder(r.Body).Decode(&payload)
  if err != nil {
    app.logger.Print(err)
  }

  //we need to get the user
  user, err := app.models.DB.GetUser(payload.Username)
  if err != nil {
    app.logger.Println("User does not exist")
    return
  } 

  hashPassword := user.Password

  err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(payload.Password))
  // Handle the error for hasing and comparing
  if err != nil {
    log.Println(err)
    _message := JSONMessage{
      Message: "Unauthorized",
    }

    js, err := json.MarshalIndent(_message, "", "\t")
    if err != nil {
      app.logger.Println(err)
    }

    w.Header().Set("Context-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(js)
    return
  }

  // Validating a users token
  var claims jwt.Claims
  claims.Subject = fmt.Sprint(user.ID)
  claims.Issued = jwt.NewNumericTime(time.Now())
  claims.NotBefore = jwt.NewNumericTime(time.Now())
  claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
  // supposed to be a unique domain you own
  claims.Issuer = "github.com/melkeydev"
  claims.Audiences = []string{"github.com/melkeydev"}

  jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.Jwt.Secret))
  if err != nil {
    fmt.Println(err)
    message := "Could not generate proper access"
    app.errorResponse(w, r, http.StatusInternalServerError, message)
    return
  }

  //app.writeJSON(w, http.StatusOK, string(jwtBytes), "Successfully logged in")
  _message := JSONMessage{
    Message : string(jwtBytes),
  }

  js, err := json.MarshalIndent(_message, "", "\t")
  if err != nil {
    app.logger.Println(err)
  }

  w.Header().Set("Context-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(js)

}


func (app *application) insertPayload(w http.ResponseWriter, r *http.Request) {
  var payload DBLoadPayload

  err := json.NewDecoder(r.Body).Decode(&payload)
  if err != nil {
    log.Println(err)
    return
  }

  dbload := &models.DBLoad{
    DBDataOne: payload.DBDataOne,
    DBDataTwo: payload.DBDataTwo,
    DBDataThree: payload.DBDataThree,
  }

  v := validator.New()

  if models.ValidateDBLoad(v, dbload); !v.Valid() {
    app.failedValidationResponse(w, r, v.Errors)
    return
  } 

  err = app.models.DB.InsertDBLoad(dbload)

  if err != nil {
    app.logger.Println(err)
  }

  _message := JSONMessage{
    Message : "Succesfully posted data to the DB",
  }

  js, err := json.MarshalIndent(_message, "", "\t")
  if err != nil {
    app.logger.Println(err)
  }

  w.Header().Set("Context-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(js)
}

func (app *application) deleteDBload(w http.ResponseWriter, r *http.Request) {
  id, err := app.readIDParam(r)
  if err != nil {
    app.notFoundResponse(w,r)
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

  err = app.writeJSON(w, http.StatusOK, envelope{"message":"data deleted Succesfully"}, nil)
  if err != nil {
    app.serverErrorResponse(w, r, err)
  }
}

//func (app *application) updateDataHandler(w http.ResponseWriter, r *http.Request) {
  //id, err := app.readIDParam(r)

  //if err != nil {
    //app.notFoundResponse(w,r)
    //return
  //}

  //if err != nil {
    //switch {
    //case errors.Is(err, models.ErrRecordNotFound):
      //app.notFoundResponse(w,r)
    //default:
      //app.serverErrorResponse(w, r, err)
    //}
    //return
  //}

  //var input struct {

  //}



//}






























