package routes_and_storage

import (
    "database/sql"
    "fmt"
    "net/http"

    "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/types"
    
)


func getCoffeeShop(w http.ResponseWriter, r *http.Request, ShopID int) {
  var shop types.CoffeeShop
  err := db.QueryRow("SELECT id, FROM coffee_shops WHERE id = ?", ShopID).Scan(&shop.ID, &shop.name)
  if err != nil {
    if err == sql.ErrNoRows {
      http.Error(w, "Coffee shop not found!", http.StatusNotFound)

    } else {
    http.Error(w, "Error retrieving coffee shop", http.StatusInternalServerError)
  }
  return
}

rows, err := db.Query("SELECT shop_id, content FROM comments WHERE shop_id = ?", ShopID)
if err != nil {
  http.Error(w, "Error retrieving comments", http.StatusInternalServerError)
  return
}
defer rows.Close()

for rows.Next() {
  var comment types.Comment
  err := rows.Scan(&comment.ID, &comment.UserID, &comment.Content)
  if err != nil {
    http.Error(w, "Error scaning comments", http.StatusInternalServerError)
    return
  }
  shop.Comments = append(shop.Comments, comment)
}

rows, err = db.Query("SELECT id, user_id, ambiance, coffee, overall FROM ratings WHERE shop_id = ?", ShopID)
if err != nil {
  http.Error(w, "Error retrieving ratings", http.StatusInternalServerError)
  return
}
defer rows.Close()

for rows.Next() {
  var rating types.Rating 
  err := rows.Scan(&rating.ID, &rating.UserID, &rating.Value)
  if err != nil {
    http.Error(w, "Error scanning ratings", http.StatusInternalServerError)
    return
  }
  shop.Ratings = append(shop.Ratings, rating)

  json.NewEncoder(w).Encode(json)

  }
}

func registerUser(w http.ResponseWriter, r *http.Request) (int, error) {
  var user types.User
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return 0, err
  }
  
  result, err :=  db.Exec("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
    user.FirstName, user.LastName, user.Email, user.Password)
  if err != nil {
    return 0, fmt.Errorf(w, "Error registering user", http.StatusInternalServerError)
  }

  id, _ := result.LastInsertId()
  if err != nil {
    return 0, fmt.Errorf(w, "Error retrieving user ID", http.StatusInternalServerError)
  }

  user.ID = int(id)
  w.Header().Set("Content-Type", "application/json") 
  w.WriteHeader(http.StatusCreated)
  if err := json.NewEncoder(w).Encode(user); err != nil {
    return 0, fmt.Errorf("error encoding user response: %w", err)
  }
  
  return int(id), nil
}

func getUser(w http.ResponseWriter, r *http.Request, UserID int) (*types.User, error) {
  var user types.User
  err := db.QueryRow("SELECT id, first_name, last_name, email FROM users WHERE id = ?", UserID).Scan(
    &user.ID, &user.FirstName, &user.LastName, &user.Email)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, fmt.Errorf("User not found", err)
      }
      return nil, fmt.Errorf("error retrieving user, %w", err)
  }
  return &user, nil
}

func addComment(w http.ResponseWriter, r *http.Request, ShopID int) (*types.Comment, error) {
  var comment types.Comment
  err := json.NewDecoder(r.Body).Decode(&comment)
  if err != nil {
    return nil, fmt.Errorf("Error with json decoding", err)
  }

  comment.ShopID = ShopID
  result, err := db.Exec("INSERT INTO comments (shop_id, content) VALUES (?, ?)",
    UserID, shopID, comment.Content)
  if err != nil {
    http.Error(w, "Error adding comment", http.StatusInternalServerError)
    return
  }

  id, _ := result.LastInsertId()
  comment.ID = int(id)

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(comment)
}

func addRating(w http.ResponseWriter, r *http.Request, ShopID int) (*types.Rating, error) {
  var rating types.Rating
  err := json.NewDecoder(r.Body).Decode(&rating)
  if err != nil {
    return nil, fmt.Errorf("Error with json decoding", err)
  }
 
  comment.ShopID = ShopID
  result, err := db.Exec("INSERT INTO ratings (shop_id, ambiance, coffee, overall) VALUES ( ?, ?, ?, ?)",
    ShopID, rating.Ambiance, rating.Coffee, rating.Overall)
  if err != nil {
    http.Error(w, "Error adding rating", http.StatusInternalServerError)
    return
  }

  id, _ := result.LastInsertId()
  rating.ID = int(id)

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(rating)

}
