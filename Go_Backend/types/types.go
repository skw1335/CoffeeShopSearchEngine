package types



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
  Ambiance  int  `json:"Ambiance rating"`
  Coffee    int  `json:"Coffee rating"`
  Overall   int  `json:"Overall rating"`
}

