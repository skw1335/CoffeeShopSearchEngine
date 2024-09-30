package store

import (
  "database/sql"
  "context"
  "time"
  "errors"
)

type Comment struct {
  ID        int64     `json:"id"`
  Content   string    `json:"content"`
  Title     string    `json:"title"`
  UserID    int64     `json:"user_id"`
  ShopID    int64     `json:"shop_id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

type CommentsStore struct {
  db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
  query := `INSERT INTO comments (content, title, user_id, shop_id)
            VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
  `
  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  err := s.db.QueryRowContext(
      ctx,
      query,
      comment.Content,
      comment.Title,
      comment.UserID,
      comment.ShopID,
  ).Scan(
        &comment.ID,
        &comment.CreatedAt,
        &comment.UpdatedAt,
  )
  if err != nil {
          return err
  }

  return nil
}

func (s *CommentsStore) GetByID(ctx context.Context, id int64) (*Comment, error) {
  query := `SELECT id, user_id, shop_id, content, title, created_at, updated_at
             FROM comments
             WHERE id = $1
             `
  var comment Comment
  err := s.db.QueryRowContext(ctx, query, id).Scan(
       &comment.ID,
       &comment.UserID,
       &comment.ShopID,
       &comment.Title,
       &comment.Content,
       &comment.CreatedAt,
       &comment.UpdatedAt,
     )
  if err != nil {
      switch {
         case errors.Is(err, sql.ErrNoRows):
          return nil, ErrNotFound
      default: 
        return nil, err
     }
    }

    return &comment, nil
   
}
