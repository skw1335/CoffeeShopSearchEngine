package db

import (
  "database/sql"
  "log"
  _ "github.com/go-sql-driver/mysql"

var db *sql.DB
var data []byte = []byte(`{}`)

// Init initializes the database connection
//
func Init(connectionString string) error {
  var err error
  db, err = sql.Open("mysql", connectionString)
  if err != nil {
    return err
  }


err = db.Ping()
  if err != nil {
    return err
}

  return initDB()
}

func Close() {
  db.Close()
}

// initDB initializes the database and creates tables if they don't exist.
func initDB() error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
            firstName VARCHAR(255) NOT NULL,
            lastName  VARCHAR(255) NOT NULL,
            email     VARCHAR(255) NOT NULL,
            password  VARCHAR(255) NOT NULL,
            createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        );

        CREATE TABLE IF NOT EXISTS coffee_shops (
            id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(100),
            address VARCHAR(255),
            
        );

        -- We'll create the comments table after we load the coffee shop data.
    `)

    // Load coffee shops from data.json
    var coffeeShops []CoffeeShop
    err = json.Unmarshal(data, &coffeeShops)
    if err != nil {
        return err
    }

    // Create the comments table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS comments (
            id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
            user_id INT,
            shop_id INT,
            content TEXT,
            FOREIGN KEY (user_id) REFERENCES users (id),
            FOREIGN KEY (shop_id) REFERENCES coffee_shops (id)
        );
    `)

    // Create the ratings table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS ratings (
            id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
            user_id INT,
            shop_id INT,
            coffee_value   INT,
            ambiance_value INT,
            overall_value  INT, 
            FOREIGN KEY (user_id) REFERENCES users (id),
            FOREIGN KEY (shop_id) REFERENCES coffee_shops (id)
        );
    `)

    // Insert coffee shops into the database
    for _, shop := range coffeeShops {
        _, err = db.Exec(`
            INSERT INTO coffee_shops (name) VALUES (?);
        `, shop.Name, shop.Address )
    }

    return err
}