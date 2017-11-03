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

// Like followers's feed
func (a *App) LikeFeedFollowers(hasNew bool) {
	if hasNew == true {
		a.getFollowers()
	}
	a.likeFollowersPosts()
}

// Like followings's feed
func (a *App) LikeFeedFollowings(hasNew bool) {
	if hasNew == true {
		a.getFollowings()
	}
	a.likeFollowingsPosts()
}

func (a *App) ShadowUser(username string, hasNew bool) {
	if hasNew == true {
		user := a.GetUserByUsername(username)
		a.getUserFollowers(user)
		a.checkUserFollowers()
	}

	a.likeAndFollow()
}
