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
  http.HandleFunc("/coffee-shop", handleCoffeeShop)
  http.HandleFunc("/user/", handleUser)

  fmt.PrintLn("Server is running on :8080")
  log.Fatal(http.ListenAndServer(":8080", nil))
}

