package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"crypto/sha256"
)

func (m *DBModel) Insert(user *User) error {
	query := `INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated}
	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail 
		default:
			return err
		}
	}
	return nil
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

func (m *DBModel) UpdateUser(user *User) error {
	query := `UPDATE users SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1 WHERE id = $5 AND version = $6 RETURNING version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *DBModel) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated, version FROM users WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
} 

// This updates the database info - not the user
func (m *DBModel) Update(load *DBLoad) error {
	// This will handle DB update race condition
	query := `UPDATE dbload SET dbdataone = $1, dbdatatwo = $2, dbdatathree = $3, version = version + 1 where id = $4 and VERSION = $5 RETURNING version`

	args := []interface{}{
		load.DBDataOne,
		load.DBDataTwo,
		load.DBDataThree,
		load.ID,
		load.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&load.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *DBModel) GetAll(DBDataOne string, filters Filters) ([]*DBLoad, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(), dbdataone, dbdatatwo, dbdatathree, id, version FROM dbload WHERE(to_tsvector('simple', dbdataone) @@ plainto_tsquery('simple', $1) OR $1='') ORDER BY %s %s, id ASC LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{DBDataOne, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	DBdata := []*DBLoad{}

	for rows.Next() {
		var data DBLoad

		err := rows.Scan(
			&totalRecords,
			&data.DBDataOne,
			&data.DBDataTwo,
			&data.DBDataThree,
			&data.ID,
			&data.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		DBdata = append(DBdata, &data)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := createMetadata(totalRecords, filters.Page, filters.PageSize)

	return DBdata, metadata, nil
}

func (m *DBModel) GetForToken(tokenScope, TokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(TokenPlaintext))

	query := `SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version FROM users INNER JOIN tokens ON users.id = tokens.user_id WHERE tokens.hash = $1 AND tokens.scope = $2 AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
