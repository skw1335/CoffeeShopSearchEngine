package main

import (
  "fmt"
  "log"
  "net/http"
  "os"

  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/api"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/db"
)

var (
	coffeeShops    = make(map[int]*CoffeeShop)
	users          = make(map[int]*User)
	mu             sync.RWMutex
	nextUserID     = 1
	nextCoffeeShopID = 1
)

func main() {
  // load env variables
  dbUser := os.Getenv("DB_USER")
  dbPass := os.Getenv("DB_PASS")
  dbHost := os.Getenv("DB_HOST")
  dbName := os.Getenv("DB_NAME")

  // initalize the database
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)
  err := db.Init(connectionString)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // set up the routes
  http.HandleFunc("/coffee-shop", api.handleCoffeeShop)
  http.HandleFunc("/user/", api.handleUser)

  fmt.PrintLn("Server is running on :8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

