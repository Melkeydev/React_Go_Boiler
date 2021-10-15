package main

import (
  "net/http"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
    next.ServeHTTP(w, r)
  })
}

