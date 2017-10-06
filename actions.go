package main

import (
	"log"
)

// Unfollow jerks
func (a *App) Unfollow() {
	// Collect data.
	log.Println("Beginning the data collection process...")
	a.getFollowers()
	a.getFollowings()

	// Compare data.
	log.Println("\nComparing the lists 🔍\n")
	a.compareLists()
	a.showList()

	a.unfollowLeeches()
}

// Like user's feed
func (a *App) LikeFeedFollowers() {
	a.getFollowers()
	a.likeFollowersPosts()
}

// Like user's feed
func (a *App) LikeFeedFollowings() {
	a.getFollowings()
	a.likeFollowingsPosts()
}
