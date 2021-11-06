package main

import (
  "io"
  "fmt"
  "log"
  "time"
  "context"
  "strings"
  "errors"
  "net/url"
  "net/http"
  "strconv"
  "encoding/json"
  "database/sql"
  "backend/types"
  "backend/validator"
  "github.com/julienschmidt/httprouter"
)

//TODO: Add to the Readme
// This will just hold useful functions that server a specific purpose for app handling

type envelope map[string]interface{}

func connectDB(cfg types.Config) (*sql.DB, error) {
  db, err := sql.Open("postgres", cfg.Db.Dsn)
  if err != nil {
    log.Fatal("unable to connect to the database")
    return nil, err
  }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

  err = db.PingContext(ctx)
  if err != nil {
    return nil, err
  }
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

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
  s := qs.Get(key)

  if s == "" {
    return defaultValue
  }

  return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string{
  csv := qs.Get(key)

  if csv == "" {
    return defaultValue
  }

  return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
  s := qs.Get(key)

  if s == "" {
    return defaultValue
  }

  i, err := strconv.Atoi(s)
  if err != nil {
    v.AddError(key, "must be an integer value")
  } 

  return i
}




























