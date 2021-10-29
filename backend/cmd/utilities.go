package main

import (
  "io"
  "fmt"
  "log"
  "strings"
  "errors"
  "net/http"
  "context"
  "strconv"
  "encoding/json"
  "database/sql"
  "backend/types"
  "backend/models"
  "github.com/julienschmidt/httprouter"
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

func (app *application) readIDParam(r *http.Request) (int64, error) {
  params := httprouter.ParamsFromContext(r.Context())

  id, err := strconv.ParseInt(params.ByName("id"),10,64)
  if err != nil {
    return 0, errors.New("invalid id parameter")
  }
  return id, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
  
  // Adds a maximum byte size to the load request
  maxBytes := 1_048_576
  r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

  // Decoder for our json payload
  dec := json.NewDecoder(r.Body)
  dec.DisallowUnknownFields()
  err := dec.Decode(dst)

  if err != nil {
    var syntaxError *json.SyntaxError
    var unmarshallTypeError *json.UnmarshalTypeError
    var invalidMarshallError *json.InvalidUnmarshalError

    switch {
    case errors.As(err, &syntaxError):
      return fmt.Errorf("body contains badly formed JSON characters %d", syntaxError.Offset)

    case errors.Is(err, io.ErrUnexpectedEOF):
      return errors.New("Badly formatted JSON in body request")

    case errors.As(err, &unmarshallTypeError):
      if unmarshallTypeError.Field != "" {
        return fmt.Errorf("body contains incorrect JSON type for field %d", unmarshallTypeError.Field)
      }
      return fmt.Errorf("Body contains incorrect JSON")

    // if there is something in our body thats empty
    case errors.Is(err, io.EOF):
      return errors.New("body must not be empty")

    case strings.HasPrefix(err.Error(), "json: unknown field"):
      // this handles when the body has an incorrect key
      fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
      return fmt.Errorf("Body contains unknown key %s", fieldName)

    case err.Error() == "http: request body too large":
      return fmt.Errorf("body must not be larger than max size")

    case errors.As(err, &invalidMarshallError):
      panic(err)

    default:
      return err
    }
  }

  err = dec.Decode(&struct{}{})
  if err != io.EOF {
    return errors.New("body must contan only valid JSON values")
  }

  return nil
}

