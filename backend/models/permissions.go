package models

import (
	"context"
	"time"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

func (m *DBModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT permissions.code	
		FROM permissions 
		INNER JOIN users_permissions ON users_permissions.permissions_id = permissions_id
		INNER JOIN users ON users_permissions.user_id = users.id
		WHERE users.id = $1
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	// every permission needs to be appended to permissions
	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil 
} 
