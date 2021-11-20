package models

import (
	"backend/validator"
	"database/sql"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	DB DBModel
}

type DBModel struct {
	DB *sql.DB
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// A generic user structure
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "password must be atleast 8 chars long")
	v.Check(len(password) <= 72, "password", "password must not  be more than 72 chars long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 72, "name", "must not be longer than 72")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password for hash use")
	}
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
