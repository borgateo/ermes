package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// New creates a new app.
func New() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	waitingTimeString := os.Getenv("WAITING_TIME")
	waitingTime, err := strconv.Atoi(waitingTimeString)
	if err != nil {
		waitingTime = 15
	}

	return &App{
		Wait:       waitingTime,
		username:   os.Getenv("USERNAME"),
		password:   os.Getenv("PASSWORD"),
		followings: map[string]bool{},
		followers:  map[string]bool{},
		leeches:    []string{},
	}
}

func main() {
	fmt.Printf("Welcome to Ermes! 👻 \n")

	app := New()
	app.Login()
	defer app.Logout()

	app.InitDB()

	//app.Unfollow()

	//app.FollowVIPFollowers("vida_nomade")

	app.LikeFeed()
}
