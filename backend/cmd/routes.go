package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
  router := httprouter.New()
  router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
  return router
}
