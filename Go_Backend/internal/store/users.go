package store

import (
  "database/sql"
  "context"
  "time"
  "errors"
  
  "golang.org/x/crypto/bcrypt"
)

var (
  ErrDuplicateEmail     = errors.New("Error: A user with that email already exists")
  ErrDuplicateUsername  = errors.New("Error: A user with that username already exists") 
)
type User struct {
  ID        int64     `json:"id"`
  Username  string    `json:"username"`
  FirstName string    `json:"first_name"`
  LastName  string    `json:"last_name"`
  Email     string    `json:"email"`
  Password  password    `json:"-"`
  CreatedAt time.Time `json:"created_at"`
}
type UsersStore struct {
  db *sql.DB
}

type password struct {
  text *string
  hash []byte
}

func (p *password) Set(text string) error {
  hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
  if err != nil {
    return err
  }

  p.text = &text
  p.hash = hash

  return nil
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
  query := `INSERT INTO users (username, first_name, last_name, email, password)
            VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
  `

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  err := tx.QueryRowContext(
    ctx,
      query,
      user.Username,
      user.FirstName,
      user.LastName,
      user.Email,
      user.Password.hash,
    ).Scan(
      &user.ID,
      &user.CreatedAt,
    )
    if err != nil {
      switch {
        case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
          return ErrDuplicateEmail
        case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
          return ErrDuplicateUsername
        default:
          return err
      }
    }

    return nil
  
}

func (s *UsersStore) GetByID(ctx context.Context, id int64)  (*User, error) {
  query := `SELECT id, username, email, password, created_at
            FROM users
            WHERE id = $1
  `  

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration) 
  defer cancel()

  user := &User{}
  err := s.db.QueryRowContext(
    ctx,
    query,
    id,
  ).Scan(
    &user.ID,
    &user.Username,
    &user.Email,
    &user.Password,
    &user.CreatedAt,
  )
  if err != nil {
    switch err {
    case sql.ErrNoRows:
      return nil, ErrNotFound
    default:
    return nil, err
    }
  }
  return user, nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
  return withTx(s.db, ctx, func(tx *sql.Tx) error {
  
  // create teh user
  if err := s.Create(ctx, tx, user); err != nil {
    return err
  }
  // create the user invite
  if err := s.createUserInvitation(ctx, tx, token, exp, user.ID);
  err != nil {
    return err
  }

  return nil
  })
}

func (s *UsersStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
  query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  _, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
  if err != nil {
    return err
  }

  return nil
}
