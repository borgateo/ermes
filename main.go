package main

import (
	"flag"
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
		log.Fatal("ERROR: something went wrong with your .env file")
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

type Fish struct{ Name string }

func main() {
	// outstanding title!
	fmt.Printf("\n███████╗██████╗ ███╗   ███╗███████╗███████╗\n")
	fmt.Printf("██╔════╝██╔══██╗████╗ ████║██╔════╝██╔════╝\n")
	fmt.Printf("█████╗  ██████╔╝██╔████╔██║█████╗  ███████╗\n")
	fmt.Printf("██╔══╝  ██╔══██╗██║╚██╔╝██║██╔══╝  ╚════██║\n")
	fmt.Printf("███████╗██║  ██║██║ ╚═╝ ██║███████╗███████║\n")
	fmt.Printf("╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝\n\n")

	// CLI flags
	// Available commands:
	// ./ermes -unfollow
	// ./ermes -followers -reset
	// ./ermes -followers -reset=true
	// ./ermes -user=username -reset=true
	unfollowPtr := flag.Bool("unfollow", false, "Unfollow the ingrates.")
	followersPtr := flag.Bool("followers", false, "Like user's followers.")
	followingsPtr := flag.Bool("followings", false, "Like user's followings.")
	resetPtr := flag.Bool("reset", false, "Fetch user's connections, and resets the DB.")
	userPtr := flag.String("user", "empty", "Follow vip's followers")

	flag.Parse()

	app := New()
	app.Login()
	defer app.Logout()

	app.InitDB()

	if *unfollowPtr == true {
		app.Unfollow()
	}

	if *userPtr != "empty" {
		app.ShadowUser(*userPtr, *resetPtr)
	}

	if *followersPtr == true {
		app.LikeFeedFollowers(*resetPtr)
	}

	if *followingsPtr == true {
		app.LikeFeedFollowings(*resetPtr)
	}

}
