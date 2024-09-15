package types

import (
  "time"
)


type User struct {
  ID        int       `json:"id"`
  firstName string    `json:"firstName"`
  lastName  string    `json:"lastName"`
  email     string    `json:"email"`
  password  string    `json:"-"`
  createdAt time.Time `json:"createdAt"`
}

type CoffeeShop struct {
  ID         int       `json:"id"`
  Name       string    `json:"name"`
  Rating     string    `json:"name"`
  Reviews    string    `json:"reviews"`
  Address    string    `json:"Address"`
  latitude   float     `json:"latitude"`
  longitude  float     `json:"longitude"`
  Comments   []Comment `json:"comments"`
  Ratings    []Rating  `json:"ratings"`
}

type Comment struct {
  ID         int       `json:"id"`
  UserID     int       `json:"userId"`
  ShopID     string    `json:"shopId"`
  Content    string    `json:"content"`
  CreatedAt  time.Time `json:"created_at"`
}

type Rating struct {
  ID        int       `json:"id"`
  UserID    int       `json:"userId"`
  ShopID    int       `json:"shopId`
  Ambiance  string    `json:"Ambiance_rating"`
  Coffee    string    `json:"Coffee_rating"`
  Overall   string    `json:"Overall_rating"`
  CreatedAt time.Time `json:"created_at"`
}

