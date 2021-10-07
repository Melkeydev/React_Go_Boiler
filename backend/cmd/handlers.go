package main

import (
  "net/http"
  "encoding/json"
)

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
