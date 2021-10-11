package main

import (
  "log"
  "net/http"
  "context"
  "encoding/json"
  "database/sql"
  "backend/types"
  "backend/models"
)

//TODO: Add to the Readme
// This will just hold useful functions that server a specific purpose for app handling

type envelope map[string]interface{}

func connectDB(ctx context.Context, cfg types.Config) (*sql.DB, error) {
  db, err := sql.Open("postgres", cfg.Db.Dsn)
  if err != nil {
    log.Fatal("unable to connect to the database")
    return nil, err
  }

  // This will create two new tables if they do not exist
  (&models.DBModel{DB:db}).CreateUsersTable(ctx)
  (&models.DBModel{DB:db}).CreateDBLoadTable(ctx)

  return db, nil
}


func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
  js, err := json.MarshalIndent(data, "", "\t")
  if err != nil {
    return err
  }

  js = append(js, '\n')

  for k, v := range headers {
    w.Header()[k] = v
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  w.Write(js)
  return nil
}


