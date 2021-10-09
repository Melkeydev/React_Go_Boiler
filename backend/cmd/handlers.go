package main

import (
  "log"
  "net/http"
  "encoding/json"
  "backend/models"
  "golang.org/x/crypto/bcrypt"
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

  _message := JSONMessage{
    Message: "Succesfully logged in",
  }

  js, err := json.MarshalIndent(_message, "", "\t")
  if err != nil {
    app.logger.Println(err)
  }

  w.Header().Set("Context-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(js)
}

















