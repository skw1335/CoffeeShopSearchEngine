package main

import (
  "log"

  mysqlCfg "github.com/go-sql-driver/mysql"
  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/mysql"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/db"
  "github.com/skw1335/CoffeeShopSearchEngine/Go_Backend/configs"
_ "github.com/golang-migrate/migrate/v4/source/file"
) 
func main() {
		cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

 driver, err := mysql.WithInstance(db, &mysql.Config{})

 m, err := migrate.NewWithDatabaseInstance(
   "file://migrate/migrations",
   "mysql"
   driver,
 )

 if err != nil {
   log.Fatal(err)
 }

 cmd := os.Args[(len(os.Args) - 1)]
 if cmd == "up" {
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
      log.fatal(err)
    }
 }
 
 if cmd == "down" {
   if err := m.Down(); err != nil && err != migrate.ErrNoChange {
      log.fatal(err)
    }
 }


