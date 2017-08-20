package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Unfollow leeches
func (a *App) Unfollow() {
	// Collect data.
	log.Println("Beginning the data collection process...")
	a.getFollowers()
	a.getFollowings()

	// Compare data.
	log.Println("Comparing the lists to each other...")
	a.compareLists()
	a.showList()

	a.unfollowLeeches()
}

// getFollowers gets users that follow us.
func (a *App) getFollowers() {
	log.Println("Collecting your 'Followers' list")
	resp, err := a.api.SelfTotalUserFollowers()
	if err != nil {
		panic(err)
	}

	for _, user := range resp.Users {
		username := strings.ToLower(user.Username)
		a.followers[username] = true
	}
}

// getFollowing gets users that we follow.
func (a *App) getFollowings() {
	log.Println("Collecting your 'Followings' list")
	resp, err := a.api.SelfTotalUserFollowing()
	if err != nil {
		panic(err)
	}

	for _, user := range resp.Users {
		username := strings.ToLower(user.Username)
		a.followings[username] = true
	}
}

// unfollow leeches -- TODO pass array to unfollow
func (a *App) unfollowLeeches() {
	var (
		counter   = 1
		remaining = len(a.leeches)
	)

	if remaining == 0 {
		fmt.Printf("No leeches - Nothing to do\n")
		return
	}

	fmt.Printf("\n 🔥 🔥 🔥  Beginning Unfollow 🔥 🔥 🔥 \n\n")

	for _, username := range a.leeches {
		if _, ok := a.followings[username]; !ok {
			fmt.Printf("[ERROR] Username %s not found in following map\n", username)
			continue
		}

		// Unfollow.
		userIDStr := a.getUserId(username)
		randomInt := random(a.Wait, a.Wait+10)
		log.Printf("- [%d of %d]: %s (UID %s) ⏰ %ds\n", counter, remaining, username, userIDStr, randomInt)

		// Convert the user ID from a string to an int.
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			panic(err)
		}

		_, err = a.api.UnFollow(int64(userID))
		if err != nil {
			log.Panicf("Got error when unfollowing %s: %s", username, err)
		}

		counter++
		time.Sleep(time.Duration(randomInt) * time.Second)
	}
}

func (a *App) showList() {
	// Sort the lists
	followings := a.sortKeys(a.followings)
	followers := a.sortKeys(a.followers)
	leeches := a.leeches

	// Sum up the numbers.
	var (
		numFollowings = len(followings)
		numFollowers  = len(followers)
		numLeeches    = len(leeches)
	)

	fmt.Printf("You've got %d followings, %d followers and leeches %d\n", numFollowings, numFollowers, numLeeches)

	for i := 0; i < len(leeches); i++ {
		fmt.Printf("- %s \n", leeches[i])
	}

}
