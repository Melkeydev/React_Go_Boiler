package models

import (
  "log"
  "time"
  "errors"
  "context"
  "database/sql"
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

func (m *DBModel) InsertDBLoad(load *DBLoad) error {
  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  query := `insert into dbload(dbdataone, dbdatatwo, dbdatathree, version) VALUES($1, $2, $3, $4)`

  _, err := m.DB.ExecContext(ctx, query, load.DBDataOne, load.DBDataTwo, load.DBDataThree)

  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func (m *DBModel) Delete(id int64) error {
  if id < 1 {
    return ErrRecordNotFound
  }

  query := `DELETE FROM dbload where id = $1`

  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  results, err := m.DB.ExecContext(ctx, query, id)
  if err != nil {
    return err
  }

  rowsAffected, err := results.RowsAffected()
  if err != nil {
    return err
  }

  if rowsAffected == 0 {
    return ErrRecordNotFound
  }

  return nil
}

func (m *DBModel) Update(load *DBLoad) error {
  // This will handle DB update race condition
  query := `UPDATE dbload SET dbdataone = $1, dbdatatwo = $2, dbdatathree = $3`

  args := []interface{} {
    load.DBDataOne,
    load.DBDataTwo,
    load.DBDataThree,
    load.Version,
  }

  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  err := m.DB.QueryRowContext(ctx, query, args...).Scan(&load.Version)
  if err != nil {
    switch{
    case errors.Is(err, sql.ErrNoRows):
      return ErrEditConflict
    default:
      return err
    }
  }

  return nil
}





























