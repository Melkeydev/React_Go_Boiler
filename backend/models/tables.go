package models

import (
  "log"
  "context"
)

// This will define how to make tables in your DB

// This is going to define out Users table
func (db *DBModel) CreateUsersTable(ctx context.Context) error {
  
  query := `create table if not exists users (
    ID SERIAL PRIMARY KEY NOT NULL,
    USERNAME TEXT NOT NULL,
    PASSWORD TEXT NOT NULL
  )`

  _, err := db.DB.ExecContext(ctx, query)

  if err != nil {
    return err
  }
  
  log.Println("Created the Users table")
  return nil

}

// This is going to define a generic db-load table
func (db *DBModel) CreateDBLoadTable(ctx context.Context) error {

  query := `create table if not exists dbload (
    ID SERIAL PRIMARY KEY NOT NULL,
    DBDATAONE TEXT NOT NULL,
    DBDATATWO TEXT NOT NULL,
    DBDATATHREE TEXT NOT NULL
  )`

  _, err := db.DB.ExecContext(ctx, query)

  if err != nil {
    return err
  }
  
  log.Println("Created the DBLoad table")
  return nil
}
