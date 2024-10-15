package store

import (
  "database/sql"
  "context"
  "time"
)

type Rating struct {
  ID        int64     `json:"id"`
  UserID    int64     `json:"user_id"`
  ShopID    int64     `json:"shop_id`
  Ambiance  int       `json:"ambiance"`
  Coffee    int       `json:"coffee"`
  Overall   int       `json:"overall"`
  Version   int       `json:"version"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

type RatingsStore struct {
  db *sql.DB
}

func (s *RatingsStore) Create(ctx context.Context, rating *Rating) error {
  query := `INSERT INTO ratings (coffee, ambiance, overall, user_id, shop_id)
            VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
  `
  err := s.db.QueryRowContext(
    ctx,
      query,
      rating.Coffee,
      rating.Ambiance,
      rating.Overall,
      rating.UserID,
      rating.ShopID,
    ).Scan(
      &rating.ID,
      &rating.CreatedAt,
  ) 
  if err != nil {
    return err
  }

  return nil
}
func (s *RatingsStore) GetByID(ctx context.Context, id int64)  (*Rating, error) {
  query := `SELECT id, rating_id, shop_id, coffee, ambiance, overall, created_at
            FROM ratings 
            WHERE id = $1
  `  

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration) 
  defer cancel()

  var rating Rating 
  err := s.db.QueryRowContext(ctx,query,id,).Scan(
    &rating.ID,
    &rating.Coffee,
    &rating.Ambiance,
    &rating.Overall,
    &rating.CreatedAt,
  )
  if err != nil {
    switch err {
    case sql.ErrNoRows:
      return nil, ErrNotFound
    default:
    return nil, err
    }
  }
  return &rating, nil
}

func (s *RatingsStore) Delete(ctx context.Context, id int64) error {
  query := `DELETE FROM ratings WHERE id = $1`

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

func (s *RatingsStore) Update(ctx context.Context, rating *Rating) error {
  query := `UPDATE ratings
            SET coffee = $1, ambiance = $2, overall = $3, version = version + 1
            WHERE id = $4 AND version = $5
            RETURNING version
  `
  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  _, err := s.db.ExecContext(ctx, query, rating.Coffee, rating.Ambiance, rating.Overall, rating.ID, rating.Version)
  if err != nil {
    return err
  }

  return nil
}
