package routes_and_storage

import (
	"database/sql"
	"fmt"

	"github.com/sikozonpc/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.firstName, user.lastName, user.Email, user.Password)
	if err != nil {
		return errVALUES ($1, $2, $3, $4)
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetCoffeeShopByID(ShopID int) (*types.CoffeeShop, error) {
  query := "SELECT ID, Name, Review, Ratings, from coffee_shops WHERE id = ?"
  row   := s.db.QueryRow(query, ShopID)

  var p types.CoffeeShop
  err := row.Scan(&p.ID, &p.Name, &p.Review, &p.Ratings)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, fmt.Errorf("no shop found with id %d", ShowID)
    }
    return nil, err
  }

  return &p, nil
}


func (s *Store) addComment (ShopID, UserID int, content string) error {
  query := `
        INSERT INTO comments (ShopID, UserID, content, created_at)
        VALUES (?, ?, ?, ?)
  `

  result, err := s.db.Exec(query, ShopID, UserID, content, time.Now())
  if err != nil {
    return fmt.Errorf("failed to add comment: %v", err)
  }

  commentID, err := result.LastInsertId()
  if err != nil {
    return fmt.Errorf("failed to get last insert ID: %v", err)
  }

  log.Printf("Added comment with ID %d to shop %d", commentID, shopID)
}

func (s *Store) addRating (ShopID, UserID int, ambiance, coffee, overall string) error {
  query := `
        INSERT INTO ratings (ShopID, UserID, ambiance, coffee, overall)
        VALUES (?, ?, ?, ?, ?)
  `

  result, err := s.db.Exec(query, ShopID, UserID, ambiance, coffee, overall)
  if err != nil {
    return fmt.Errorf("failed to add rating: %v", err)
  }

  ratingID, err := result.LastInsertId()
  if err != nil {
    return fmt.Erorrf("failed to get last insert ID: %v", err)
  }

  log.Printf("Added rating with ID %d to shop %d", ratingID, ShopID)
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.firstName,
		&user.lastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
