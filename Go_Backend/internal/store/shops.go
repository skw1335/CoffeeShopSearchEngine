package store

import (
  "database/sql"
  "context"
  "errors"
)

type Shop struct {
  ID        int64     `json:"id"`
  ShopName  string    `json:"shop_name"`
  Review    string    `json:"review"`
  Rating    float32   `json:"rating"`
  Address   string    `json:"address"`
  Latitude  float64   `json:"latitude"` 
  Longitude float64   `json:"longitude"`
}


type ShopStore struct {
  db *sql.DB
}


func (s *ShopStore) GetByID(ctx context.Context, id int64) (*Shop, error) {
  query := `SELECT id, shop_name, review, rating, address, latitude, longitude
             FROM coffee_shops 
             WHERE id = $1
             `

  ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
  defer cancel()

  var shop Shop 
  err := s.db.QueryRowContext(ctx, query, id).Scan(
       &shop.ID,
       &shop.ShopName,
       &shop.Review,
       &shop.Rating,
       &shop.Address,
       &shop.Longitude,
       &shop.Latitude,
     )
  if err != nil {
      switch {
         case errors.Is(err, sql.ErrNoRows):
          return nil, ErrNotFound
      default: 
        return nil, err
     }
    }

    return &shop, nil
   
}
