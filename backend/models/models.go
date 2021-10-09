package models 

import (
  "database/sql"
)

type Models struct {
  DB DBModel
}

type DBModel struct {
  DB *sql.DB
}

func NewModels(db *sql.DB) Models {
  return Models{
    DB: DBModel{DB: db},
  }
}

// A generic user structure
type User struct {
  ID int64 `json:"id"`
  Username string `json:"username"`
  Password string `json"-"`
}

// A generic payload structure got API calls
type Payload struct {
  SampleOne string `json:"sample_one"`
  SampleTwo string `json:"sample_two"`
  SampleThree string `json:"sample_three"`
} 

type DBLoad struct {
  DBDataOne string `json:db_data_one`
  DBDataTwo string `json:db_data_two`
  DBDataThree string `json:db_data_three`
}






