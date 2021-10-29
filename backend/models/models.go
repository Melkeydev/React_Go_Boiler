package models

import (
	"backend/validator"
	"database/sql"
	"errors"
)

type Models struct {
	DB DBModel
}

type DBModel struct {
	DB *sql.DB
}

var (
	ErrRecordNotFound = errors.New("Record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// A generic user structure
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json"-"`
	Version  int32  `json:version`
}

// A generic payload structure got API calls
type Payload struct {
	SampleOne   string `json:"sample_one"`
	SampleTwo   string `json:"sample_two"`
	SampleThree string `json:"sample_three"`
}

type DBLoad struct {
	DBDataOne   string `json:db_data_one`
	DBDataTwo   string `json:db_data_two`
	DBDataThree string `json:db_data_three`
	ID          int64  `json:id`
	Version     int32  `json:version`
}

func ValidateDBLoad(v *validator.Validator, dbload *DBLoad) {
	v.Check(dbload.DBDataOne != "", "dbdataone", "data for field one must be provided")
	v.Check(len(dbload.DBDataOne) <= 500, "dbdataone", "data must be less than 500 chars")
	v.Check(dbload.DBDataTwo != "", "dbdatatwo", "data for field two must be provided")
	v.Check(len(dbload.DBDataTwo) <= 500, "dbdatatwo", "data must be less than 500 chars")
	v.Check(dbload.DBDataThree != "", "dbdatathree", "data for field three must be provided")
	v.Check(len(dbload.DBDataThree) <= 500, "dbdatathree", "data must be less than 500 chars")
}
