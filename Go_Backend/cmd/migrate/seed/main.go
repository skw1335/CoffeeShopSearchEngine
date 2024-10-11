package main

import (
  "log"
  "CoffeeMap/internal/env"
  "CoffeeMap/internal/db"
  "CoffeeMap/internal/store"
  _ "github.com/lib/pq"
)

func main() {
  addr := env.GetString("DB_ADDR", "postgres://postgres:Theallure1!@localhost/coffeeMap?sslmode=disable")
  conn, err := db.New(addr, 3, 3, "15m")
  if err != nil {
    log.Fatal(err)
  }

  defer conn.Close()
  store := store.NewStorage(conn)

  db.Seed(store)
  
}
