package store

import (
	"database/sql"
	"context"
)

type Role struct {
	ID					int64  `json:"id"`
	Name				string `json:"name"`
	Description	string `json:"description"`
	Level				int64  `json:"Level"`
}

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	query := `SELECT id, name, description, level FROM roles WHERE name = $1`

	role := &Role{}
	err := s.db.QueryRowContext(ctx, query, slug).Scan(&role.ID, &role.Name, &role.Description, &role.Level)
	if err != nil {
		return nil, err
	}

	return role, err
}
