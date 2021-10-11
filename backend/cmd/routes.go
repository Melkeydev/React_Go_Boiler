package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
  router := httprouter.New()

    //Add our custom error handling 
  router.NotFound = http.HandlerFunc(app.notFoundResponse)
  router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

  router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)


  /// TODO: ALL NEW VERSIONS
  router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
  router.HandlerFunc(http.MethodPost, "/register", app.registerUser)
  router.HandlerFunc(http.MethodPost, "/login/", app.login)
  router.HandlerFunc(http.MethodPost, "/post_data/", app.insertPayload)
  return router
}
