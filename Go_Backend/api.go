package main

import (
  "encoding/json"
  "net/http"
  "strconv"
  "strings"
 
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/db"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/types"
  _ "github.com/go-sql-driver/mysql"
)


type APIServer struct {
  addr string
}

func NewApiServer(addr string) *APIServer {
  return &APIServer {
    addr: addr,
  }
}

func (s *APIServer) Run() error {
  router := http.NewServerMux()
  
  // Get User
  router.HandleFunc("GET users/{UserID}", func(w http.ResponseWriter, r *http.Request) { 
    UserID := r.PathValue("UserID")
    user, err := getUser(w, r, UserID)
    if err != nil {
        if errors.Is(err, "user not found") {
          http.Error(w, err.Error(), http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
      return
  }
  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(user); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
    }
  })
  // Register a user
  router.HandleFunc("POST /users/register", func(w http.ResponseWriter, r *http.Request) {
    userID, err := registerUser(w, r)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
})
  // Get Coffee Shop
  router.HandleFunc("GET coffee_shops/{ShopID}", func (w http.ResponseWriter, r *http.Request) {
    ShopID := r.PathValue(ShopID)
    getCoffeeShop(w, r, ShopID)
    w.Write([]byte("Shop ID: " + ShopID))
  }
  // Post Comment
  //
  router.HandleFunc("POST /coffee-shop/{ShopID}/comments", func (w http.ResponseWriter, r *http.Request) {
    ShopID := r.PathValue(ShopID)
    addComment(w, r, ShopID)
    w.Write([]byte("Shop ID: " + ShopID))
  })
  // Post Rating
  //
  router.HandleFunc("POST /coffee-shop/{ShopID}/rating")
    ShopID := r.PathValue(ShopID)
    addComment(w, r, ShopID)
    w.Write([]byte("Shop ID: " + ShopID))

  }
    server := http.Server{
    Addr: s.addr,
    Handler: router,
  }

  log.Printf("Server running on %s," s.addr)
}


func getCoffeeShop(w http.ResponseWriter, r *http.Request, ShopID int) {
  var shop types.CoffeeShop
  err := db.QueryRow("SELECT id, name FROM coffee_shops WHERE id = ?", ShopID).Scan(&shop.ID, &shop.name)
  if err != nil {
    if err == sql.ErrNoRows {
      http.Error(w, "Coffee shop not found!", http.StatusNotFound)

    } else {
    http.Error(w, "Error retrieving coffee shop", http.StatusInternalServerError)
  }
  return
}

rows, err := db.Query("SELECT id, content FROM comments WHERE shop_id = ?", ShopID)
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
  if err := json.NewEncoder(w).Encode(user) err != nil {
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

func addComment(w http.ResponseWriter, r *http.Request, ShopID, UserID int) (*types.Comment, error) {
  var comment types.Comment
  err := json.NewDecoder(r.Body).Decode(&comment)
  if err != nil {
    return nil, fmt.Errorf("Error with json decoding", err)
  }

  comment.UserID = UserID
  comment.ShopID = ShopID
  result, err := db.Exec("INSERT INTO comments (user_id, shop_id, content) VALUES (?, ?, ?)",
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

func addRating(w http.ResponseWriter, r *http.Request, UserID, ShopID int) (*types.Rating, error) {
  var rating types.Rating
  err := json.NewDecoder(r.Body).Decode(&rating)
  if err != nil {
    return nil, fmt.Errorf("Error with json decoding", err)
  }
 
  comment.UserID = UserID
  comment.ShopID = ShopID
  result, err := db.Exec("INSERT INTO ratings (user_id, shop_id, ambiance, coffee, overall) VALUES (?, ?, ?, ?, ?, ?)",
    rating.UserID, ShopID, rating.Ambiance, rating.Coffee, rating.Overall)
  if err != nil {
    http.Error(w, "Error adding rating", http.StatusInternalServerError)
    return
  }

  id, _ := result.LastInsertId()
  rating.ID = int(id)

  w.WriteHEader(http.StatusCreated)
  json.NewEncoder(w).Encode(rating)

}


