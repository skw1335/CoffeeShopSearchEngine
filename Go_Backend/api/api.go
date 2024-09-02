package api

import (
  "encoding/json"
  "net/http"
  "strconv"
  "string"

  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/db"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/types"
)

func HandleCoffeeShop(w http.ResponseWriter, r *http.Request) {
  segments := strings.Split(r.URL.Path, "/")
  if len(segments < 3) {
    http.Error(w, "Invalid URL", http.StatusBadRequest)
    return
  }

  ShopID, err: strconv.Atoi(segments[2])
  if err != nil {
    http.Error(w, "Invalid shop ID," http.StatusBadRequest)
    return
  }

  action := ""
  if len(segments) > 3 {
    action = segments[3]
  

  switch r.Method {
  case http.MethodPost:
    switch action {
    case "comment":
      addComment(w, r, ShopID)
    case "rate":
      addRating(w, r, ShopID)
    default:
      http.Error(w, "Invalid action", http.StatusBadRequest)

    }
  case http.MethodGet:
    getCoffeeShop(w, r, ShopID)
  default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
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

func handleUser(w http.ResponseWriter, r *http.Request) {
  segments := string.Split(r.URL.Path, "/")
  if len(segments) < 3 {
    http.Error(w "Invalid URL", http.StatusBadRequest)
    return
  }
  action := segments[2]

  switch  r.Method {
  case http.MethodPost:
    if action == "register"
      registerUser(w, r)
    } else {
      http.Error(w, "Invalid action", http.StatusBadRequest)
    }
  case http.MethodGet:
    UserID, err := strconv.Atoi(action)
    if err != nil {
      http.Error(w, "Invalid user ID", http.StatusBadRequest)
      return
    }
    getUser(w, r, UserID)
  default:
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
}

func registerUser(w http.ResponseWriter, r *http.Request) {
  var user types.User
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
  }
  
  result, err :=  db.Exec("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
    user.FirstName, user.LastName, user.Email, user.Password)
  if err != nil {
    http.Error(w, "Error registering user", http.StatusInternalServerError)
    return
  }

  id, _ := result.LastInsertId()
  user.ID int(id)

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(user)

}

func getUser(w http.ResponseWriter, r *http.Request, UserID int)
  var user types.User
  err := db.QueryRow("SELECT id, first_name, last_name, email FROM users WHERE id = ?", UserID).Scan(
    &user.ID, &user.FirstName, &user.LastName, &user.Email)
  if err != nil {
    if err == sql.ErrNoRows {
      http.Error(w, "User not found", http.StatusNotFound)

    } else {
      http.Error(w, "Error retrieving user", http.StatusInternalServerError)
    }
    return
  }

  json.NewEncoder(w).Encode(user)
}

func addComment(w http.ResponseWriter, r *http.Request, ShopID int) {
  var comment types.Comment
  err := json.NewDecoder(r.Body).Decode(&comment)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  result, err := db.Exec("INSERT INTO comments (user_id, shop_id, content) VALUES (?, ?, ?)",
    comment.UserID, shopID, comment.Content)
  if err != nil {
    http.Error(w, "Error adding comment", http.StatusInternalServerError)
    return
  }

  id, _ := result.LastInsertId()
  comment.ID = int(id)

  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(comment)
}

func addRating(w http.ResponseWriter, r *http.Request, ShopID int) {
  var rating types.Rating
  err := json.NewDecoder(r.Body).Decode(&rating)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  
}
  result, err := db.Exec("INSERT INTO ratings (user_id, shop_id, ambiance, coffee, overall) VALUES (?, ?, ?, ?, ?, ?)"
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
  


