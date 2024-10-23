package store

import (
  "database/sql"
  "context"
  "time"
  "errors"
  "encoding/hex"
  "crypto/sha256"
  
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
  Password  password  `json:"-"`
  CreatedAt time.Time `json:"created_at"`
  IsActive  bool      `json:"is_active"`
	RoleID		int64 		`json:"role_id"`
	Role			Role 			`json:"role"`
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
  query := `INSERT INTO users (username, first_name, last_name, email, password, role_id)
            VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at
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
			user.RoleID,
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
  query := `SELECT users.id, username, email, password, created_at, roles.*
            FROM users
						JOIN roles ON (users.role_id = roles.id)
            WHERE user.id = $1
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
    &user.Password.hash,
    &user.CreatedAt,
		&user.Role.ID,
		&user.Role.Level,
		&user.Role.Description,
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

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password, created_at
						FROM users
						WHERE email = $1 and is_active = true
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}
  err := s.db.QueryRowContext(ctx,query,email).Scan(
    &user.ID,
    &user.Username,
    &user.Email,
    &user.Password.hash,
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
func (s *UsersStore) getUserByInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
  query := `
    SELECT u.id, u.username, u.email, u.created_at, u.is_active
    FROM users u
    JOIN user_invitations ui ON u.id = ui.user_id
    WHERE ui.token = $1 AND ui.expiry > $2
  `
  hash := sha256.Sum256([]byte(token))
  hashToken := hex.EncodeToString(hash[:])

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  user := &User{}
  err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
    &user.ID,
    &user.Username,
    &user.Email,
    &user.CreatedAt,
    &user.IsActive,
  )
  if err != nil {
    switch err{
      case sql.ErrNoRows:
        return nil, ErrNotFound
      default: 
        return nil, err
    }
  }

  return user, nil
}
func (s *UsersStore) Activate(ctx context.Context, token string) error {
  // 1. find the user that the token belongs to
  // 2. update the user 
  // 3. clean the invitations
  return withTx(s.db, ctx, func(tx *sql.Tx) error {
  user, err := s.getUserByInvitation(ctx, tx, token)
  if err != nil {
    return err
  }

  user.IsActive = true
  if err := s.update(ctx, tx, user); err != nil {
    return err
  }

  if err := s.deleteUserInvitations(ctx, tx, user.ID); err != nil {
    return err
    }
  return nil
  })
}

func (s *UsersStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
  query := `UPDATE users SET username = $1,
            email = $2, is_active = $3
            WHERE id = $4
  `
  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  _, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
  if err != nil {
    return err
  }
  return nil
}

func (s *UsersStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, UserID int64) error {
  query := `DELETE FROM user_invitations WHERE user_id = $1`

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  _, err := tx.ExecContext(ctx, query, UserID)
  if err != nil {
    return err
  }

  return nil

}

func (s *UsersStore) Delete(ctx context.Context, UserID int64) error {
  return withTx(s.db, ctx, func(tx *sql.Tx) error {
    if err := s.delete(ctx, tx, UserID); err != nil {
      return err
    }
    if err := s.deleteUserInvitations(ctx, tx, UserID); err != nil {
      return err
    }
    return nil
  })
}

func (s *UsersStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
  query := `DELETE FROM users WHERE id = $1`

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  _, err := tx.ExecContext(ctx, query, id)
  if err != nil {
    return err
  }


  return nil
}

