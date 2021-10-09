package main

import (
  "log"
  "context"
  "database/sql"
  "backend/types"
  "backend/models"
)

//TODO: Add to the Readme
// This will just hold useful functions that server a specific purpose for app handling

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
