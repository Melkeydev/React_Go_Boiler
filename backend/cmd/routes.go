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

  // we need to put the authetnication wrapper on each route

  router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
  router.HandlerFunc(http.MethodGet, "/v1/status", app.statusHandler)
  router.HandlerFunc(http.MethodGet, "/v1/data/:id", app.requireActivatedUser(app.getData))
  router.HandlerFunc(http.MethodGet, "/v1/data", app.requireActivatedUser(app.listAllDBData))
  router.HandlerFunc(http.MethodPost, "/v1/register", app.registerUser)
  //router.HandlerFunc(http.MethodPost, "/v1/login/", app.login)
  router.HandlerFunc(http.MethodPost, "/v1/post_data/", app.requireActivatedUser(app.insertPayload))
  router.HandlerFunc(http.MethodPatch, "/v1/data/:id", app.requireActivatedUser(app.updateDBData))
  router.HandlerFunc(http.MethodDelete, "/v1/data/:id", app.requireActivatedUser(app.deleteDBload))
  router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUser)
  router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

  return app.recoverPanic(app.rateLimit(app.enableCORS(app.authenticate(router))))
}
