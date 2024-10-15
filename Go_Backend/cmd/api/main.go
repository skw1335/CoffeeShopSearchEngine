package main

import (
  "log"
  "os"
  "time"
  "CoffeeMap/internal/env"
  "CoffeeMap/internal/store"
  "CoffeeMap/internal/db"

  _ "github.com/lib/pq"
)

const version = "0.0.1"

//	@title			CoffeeMap API written in Golang 
//	@description	API for CoffeeMap, an interactive map for coffee shops written in Golang 
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {
  dbAddr := os.Getenv("DB_ADDR")
  cfg := config {
    addr: env.GetString("ADDR", "3030"),
    apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
    db: dbConfig{
      addr: env.GetString("DB_ADDR", dbAddr),
      maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
      maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
      maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
   },
    env: env.GetString("ENV", "development"),
    mail: mailConfig{
      exp: time.Hour * 24 * 3, //3 days
    },
  }

  db, err := db.New(
    cfg.db.addr,
    cfg.db.maxOpenConns,
    cfg.db.maxIdleConns,
    cfg.db.maxIdleTime,
  )
  if err != nil {
    log.Panic(err)
  }

  defer db.Close()
  log.Println("database connection pool established")
  
  store := store.NewStorage(db)
  app := &application{
      config: cfg,
      store: store, 
  }
 
  mux := app.mount()
  log.Fatal(app.run(mux))
}
