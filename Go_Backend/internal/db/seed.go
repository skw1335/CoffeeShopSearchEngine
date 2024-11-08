package db

import (
  "CoffeeMap/internal/store"
  "fmt"
  "log"
  "math/rand"
  "context"
  "database/sql"
)

const shoplen = 238

var usernames = []string {
 "Olivia", "Liam", "Emma", "Noah", "Sophia", "Jackson", "Ava", "Lucas", "Mia", "Aiden",
		"Amelia", "Elijah", "Isabella", "Oliver", "Aria", "James", "Ella", "Benjamin", "Charlotte", "Mason",
		"Harper", "Ethan", "Scarlett", "Alexander", "Abigail", "Henry", "Emily", "Sebastian", "Madison", "Jack",
		"Lily", "Owen", "Avery", "Leo", "Sofia", "Jacob", "Zoe", "William", "Layla", "Daniel", "Grace",
		"Matthew", "Riley", "Carter", "Eleanor", "Luke", "Hannah", "Wyatt", "Nora", "Gabriel", "Hazel",
		"David", "Ellie", "Isaac", "Luna", "Jayden", "Victoria", "Logan", "Penelope", "Joseph", "Aurora",
		"Samuel", "Brooklyn", "Levi", "Savannah", "Dylan", "Bella", "Anthony", "Aubrey", "Andrew", "Claire",
		"Isaiah", "Willow", "Christopher", "Stella", "Joshua", "Lucy", "Julian", "Violet", "Nathan", "Paisley",
		"Ryan", "Alice", "Grayson", "Anna", "Caleb", "Natalie", "Christian", "Skylar", "Hunter", "Leah",
		"Adam", "Serenity", "Jonathan", "Eva", "Charles", "Mila", "Thomas", "Chloe", "Aaron", "Evelyn",
}

var contents = []string {
    "Great atmosphere and friendly staff!",
	"Best coffee in town, hands down.",
	"Cozy spot but a bit pricey.",
	"Loved the pastries, but the coffee was average.",
	"Perfect place to work with free WiFi.",
	"Service was slow but the coffee made up for it.",
	"Overhyped and underwhelming.",
	"Quiet and perfect for reading a book.",
	"Excellent cold brew, but the lattes were too sweet.",
	"Not bad, but I’ve had better elsewhere.",
	"Always busy but worth the wait.",
	"Baristas are very knowledgeable and friendly.",
	"The cappuccino was perfect!",
	"I didn’t like the seating arrangement.",
	"Great place, but parking is a nightmare.",
	"Decent coffee, but not a fan of the music they play.",
	"The outdoor seating is lovely on a sunny day.",
	"Amazing espresso and friendly service.",
	"Clean space and modern vibe, I’ll be back.",
	"The coffee was burnt, very disappointing.",
}

var titles = []string {
  "Fantastic Vibe with Friendly Service",
	"Best Coffee in Town!",
	"Cozy but Pricey",
	"Tasty Pastries, So-So Coffee",
	"Perfect Workspace with WiFi",
	"Slow Service, Great Coffee",
	"Overhyped Experience",
	"Ideal Quiet Spot for Reading",
	"Great Cold Brew, Overly Sweet Lattes",
	"Decent but Not Memorable",
	"Busy but Worth the Wait",
	"Friendly and Knowledgeable Baristas",
	"Perfect Cappuccino!",
	"Uncomfortable Seating",
	"Great Coffee, Bad Parking",
	"Good Coffee, Terrible Music",
	"Lovely Outdoor Seating",
	"Amazing Espresso and Friendly Vibes",
	"Modern and Clean, Will Return",
	"Burnt Coffee, Disappointing Visit", 
}


func Seed(store store.Storage, db *sql.DB) {
  ctx := context.Background()
  
  users := generateUsers(100)
  tx, _ := db.BeginTx(ctx, nil)
  for _, user := range users {
    if err := store.Users.Create(ctx, tx, user); err != nil {
      _ = tx.Rollback()
      log.Println("Error creating user:", err)
      return
    }
  }

  tx.Commit()

  comments := generateComments(200, users)
  for _, comment := range comments {
    if err := store.Comments.Create(ctx, comment); err != nil {
      log.Println("Error creating comment:", err)

    }
  }
  ratings := generateRatings(200, users)
  for _, rating := range ratings {
    if err := store.Ratings.Create(ctx, rating); err != nil {
      log.Println("Error creating rating", err)
    }
  } 
  log.Println("Seeding Complete")
}

func generateUsers(num int) []*store.User {
  users := make([]*store.User, num)


  for i := 0; i < num; i++ {
    users[i] = &store.User{
      Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
      Email: usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
      Role: store.Role{
	      Name: "user",
      },
    }
  }

  return users
}

func generateComments(num int, users []*store.User) []*store.Comment {
  comments := make([]*store.Comment, num)
  for i := 0; i < num; i++ {
    user := users[rand.Intn(len(users))]

    comments[i] = &store.Comment{
      UserID: user.ID,
      ShopID: int64(rand.Intn(shoplen)),
      Title: titles[int64(rand.Intn(len(titles)))],
      Content: titles[int64(rand.Intn(len(contents)))],
    }
  }

  return comments
}


func generateRatings(num int, users []*store.User) []*store.Rating {
  ratings := make([]*store.Rating, num)
  for i := 0; i < num; i++ {
    user := users[rand.Intn(len(users))]

    ratings[i] = &store.Rating{
      UserID:   user.ID,
      ShopID:   int64(rand.Intn(shoplen)),
      Coffee:   rand.Intn(10),
      Ambiance: rand.Intn(10),
      Overall:  rand.Intn(10),
    }
  } 
  return ratings
} 
