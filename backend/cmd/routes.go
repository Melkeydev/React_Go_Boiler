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
  router.HandlerFunc(http.MethodGet, "/v1/status", app.statusHandler)
  router.HandlerFunc(http.MethodGet, "/v1/data/:id", app.getData)
  router.HandlerFunc(http.MethodGet, "/v1/data", app.listAllDBData)
  router.HandlerFunc(http.MethodPost, "/v1/register", app.registerUser)
  router.HandlerFunc(http.MethodPost, "/v1/login/", app.login)
  router.HandlerFunc(http.MethodPost, "/v1/post_data/", app.insertPayload)
  router.HandlerFunc(http.MethodPatch, "/v1/data/:id", app.updateDBData)
  router.HandlerFunc(http.MethodDelete, "/v1/data/:id", app.deleteDBload)

  return app.recoverPanic(app.rateLimit(app.enableCORS(router)))
}
