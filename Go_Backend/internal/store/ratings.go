package store

import (
  "database/sql"
  "context"
  "time"
)

type Rating struct {
  ID        int64       `json:"id"`
  UserID    int64       `json:"user_id"`
  ShopID    int64       `json:"shop_id`
  Ambiance  int    `json:"ambiance_rating"`
  Coffee    int  `json:"coffee_rating"`
  Overall   int    `json:"overall_rating"`
  Version   int    `json:"version"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
}

type RatingsStore struct {
  db *sql.DB
}

func (s *RatingsStore) Create(ctx context.Context, rating *Rating) error {
  query := `INSERT INTO ratings (ambiance_rating, coffee_rating, overall_rating, user_id, shop_id)
            VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
  `
  err := s.db.QueryRowContext(
    ctx,
      query,
      rating.Ambiance,
      rating.Coffee,
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
