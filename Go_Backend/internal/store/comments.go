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
  Version   int       `json:"version"`
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
  query := `SELECT id, user_id, shop_id, content, title, created_at, updated_at, version
             FROM comments
             WHERE id = $1
             `

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  var comment Comment
  err := s.db.QueryRowContext(ctx, query, id).Scan(
       &comment.ID,
       &comment.UserID,
       &comment.ShopID,
       &comment.Title,
       &comment.Content,
       &comment.CreatedAt,
       &comment.UpdatedAt,
       &comment.Version,
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

func (s *CommentsStore) Delete(ctx context.Context, id int64) error {
  query := `DELETE FROM posts WHERE id = $1`

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  res, err := s.db.ExecContext(ctx, query, id)
  if err != nil {
    return err
  }

  rows, err := res.RowsAffected()
  if err != nil {
    return err
  }

  if rows == 0 {
    return ErrNotFound
  }

  return nil
} 

func (s *CommentsStore) Update(ctx context.Context, comment *Comment) error {
  query := `UPDATE comments
            SET title = $1, content = $2, version = version + 1
            WHERE id = $3 AND version = $4
            RETURNING version
  `
  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  _, err := s.db.ExecContext(ctx, query, comment.Title, comment.Content, comment.ID, comment.Version)
  if err != nil {
    return err
  }

  return nil
}
