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
  }
  Users interface {
    Create(context.Context, *User) error
  }
  Ratings interface {
    Create(context.Context, *Rating) error
  }
}

func NewStorage(db *sql.DB) Storage {
  return Storage{
    Comments: &CommentsStore{db},
    Users: &UsersStore{db},
    Ratings: &RatingsStore{db},
  }
}
