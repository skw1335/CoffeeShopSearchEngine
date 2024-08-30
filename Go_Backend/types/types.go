package types



type User struct {
  ID        int    `json:"id"`
  FirstName string `json:"firstName"`
  LastName  string `json:"lastName"`
  Email     string `json:"email"`
  Password  string `json:"-"`
}

type CoffeeShop struct {
  ID         int       `json:"id"`
  Name       string    `json:"name"`
  Comments   []Comment `json:"comments"`
  Ratings    []Ratings `json:"ratings"`
}

type Comment struct {
  UserID     int      `json:"userId"`
  ShopID     string   `json:"shopId"`
  Content    string   `json:"content"`
}

type Ratings struct {
  ID        int     `json:"id"`
  UserID    int     `json:"userId"`
  ShopID    int     `json:"shopId`
  Ambiance  string  `json:"Ambiance rating"`
  Coffee    string  `json:"Coffee rating"`
  Overall   string  `json:"Overall rating"`
}

