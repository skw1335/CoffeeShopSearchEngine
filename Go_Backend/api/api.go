package api 

import (
  "encoding/json"
  "net/http"
  "strconv"
  "strings"
  "log"
  "database/sql" 
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/db"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/types"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/routes"

)


type APIServer struct {
  addr string
  db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
  return &APIServer {
    addr: addr,
    db:   db, 
  }
}

func (s *APIServer) Run() error {
  router := http.NewServeMux()
  
  // Get User
  router.HandleFunc("GET users/{UserID}", func(w http.ResponseWriter, r *http.Request) { 
    id := r.PathValue("UserID")
    UserID, err := strconv.Atoi(id) 
      if err != nil {
        log.Fatal(err)
      }
    user, err := getUser(w, r, UserID)
    if err != nil {
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

    w.Header().Set("Content-Type", "application/json")
    response := map[string]int{"userID": userID}
    if err := json.NewEncoder(w).Encode(response); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusCreated)
})
  // Get Coffee Shop
  router.HandleFunc("GET coffee_shops/{ShopID}", func (w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("ShopID") 
    ShopID, err := strconv.Atoi(id)
      if err != nil {
        log.Fatal(err)
      }
    getCoffeeShop(w, r, ShopID)
    w.Write([]byte("Shop ID: " + id))

  })
  // Post Comment
  //
  router.HandleFunc("POST /coffee-shop/{ShopID}/comments", func (w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("ShopID")
    ShopID, err := strconv.Atoi(id)
      if err != nil {
        log.Fatal(err)
    }
    addComment(w, r, ShopID)
    w.Write([]byte("Shop ID: " + id))
  })
  // Post Rating
  //
  router.HandleFunc("POST /coffee-shop/{ShopID}/rating", func (w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("ShopID")
    ShopID, err := strconv.Atoi(id)
      if err != nil {
        log.Fatal(err)
      }
    addRating(w, r, ShopID)
    w.Write([]byte("Shop ID: " + id))
  })
  log.Printf("Server running on %s", s.addr)
  
  return http.ListenAndServe(s.addr, router)
}

