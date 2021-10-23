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

func (m *DBModel) GetData(id int64) (*DBLoad, error) {
  if id < 1 {
    return nil, ErrRecordNotFound
  }

  query := `SELECT dbdataone, dbdatatwo, dbdatathree, version from dbload where id = $1`

  var load DBLoad

  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  err := m.DB.QueryRowContext(ctx, query, id).Scan(
    &load.DBDataOne,
    &load.DBDataTwo,
    &load.DBDataThree,
    &load.Version,
  )

  if err != nil {
    switch {
    case errors.Is(err, sql.ErrNoRows):
      return nil, ErrRecordNotFound
    default:
      return nil, err
    }
  }

  return &load, nil
}

func (m *DBModel) InsertDBLoad(load *DBLoad) error {
  query := `insert into dbload(dbdataone, dbdatatwo, dbdatathree) VALUES($1, $2, $3) returning version`

  args := []interface{}{load.DBDataOne, load.DBDataTwo, load.DBDataThree}

  ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
  defer cancel()

  return m.DB.QueryRowContext(ctx, query, args...).Scan(&load.DBDataOne)
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









