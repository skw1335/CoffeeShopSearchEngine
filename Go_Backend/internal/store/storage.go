package store

import (
  "context"
  "database/sql"
  "time"
  "errors"
)

var (
  QueryTimeoutDuration = time.Second * 5
  ErrNotFound = errors.New("Resource not found")
)
type Storage struct {
  Comments interface {
    GetByID(context.Context, int64) (*Comment, error)
    Create(context.Context, *Comment) error
    Delete(context.Context, int64) error
    Update(context.Context, *Comment) error
  }
  Users interface {
    GetByID(context.Context, int64) (*User, error)
    Create(context.Context, *sql.Tx, *User) error
    CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
    Activate(context.Context, string) error
  }
  Ratings interface {
    GetByID(context.Context, int64) (*Rating, error)
    Create(context.Context, *Rating) error
    Delete(context.Context, int64) error
    Update(context.Context, *Rating) error
  }
  Shops interface {
    GetByID(context.Context, int64) (*Shop, error)
  }
}

func NewStorage(db *sql.DB) Storage {
  return Storage{
    Comments: &CommentsStore{db},
    Users: &UsersStore{db},
    Ratings: &RatingsStore{db},
  }
}


func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
  tx, err := db.BeginTx(ctx, nil)
  if err != nil {
    return nil
  }

  if err := fn(tx); err != nil {
    _ = tx.Rollback()
    return err
  }

  return tx.Commit()
}
