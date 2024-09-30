package store

import (
  "database/sql"
  "context"
  "time"
  
)

type User struct {
  ID        int64     `json:"id"`
  Username  string    `json:"username"`
  FirstName string    `json:"first_name"`
  LastName  string    `json:"last_name"`
  Email     string    `json:"email"`
  Password  string    `json:"-"`
  CreatedAt time.Time `json:"created_at"`
}
type UsersStore struct {
  db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
  query := `INSERT INTO users (username, first_name, last_name, email, password)
            VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
  `

  err := s.db.QueryRowContext(
    ctx,
      query,
      user.Username,
      user.FirstName,
      user.LastName,
      user.Email,
      user.Password,
    ).Scan(
      &user.ID,
      &user.CreatedAt,
    )
    if err != nil {
      return err
    }

    return nil
  
}
