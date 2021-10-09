package models

import (
  "log"
  "time"
  "context"
)

func (m *DBModel) RegisterUser(user User) error {
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  query := `insert into users(username, password) VALUES ($1, $2)`

  _, err := m.DB.ExecContext(ctx, query, user.Username, user.Password)
  if err != nil {
    log.Println(err)
    return err
  }
  
  return nil
}

func (m *DBModel) GetUser(username string) (*User, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  query := `select id, username, password from users where username = $1`

  row := m.DB.QueryRowContext(ctx, query, username)

  var user User

  err := row.Scan(
    &user.ID,
    &user.Username,
    &user.Password,
  )

  if err != nil {
    return nil, err
  }

  return &user, nil
}

func (m *DBModel) InsertDBLoad(load DBLoad) error {
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  query := `insert into dbload(dbdataone, dbdatatwo, dbdatathree) VALUES($1, $2, $3)`

  _, err := m.DB.ExecContext(ctx, query, load.DBDataOne, load.DBDataTwo, load.DBDataThree)

  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}





























