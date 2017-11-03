package main

import (
	"log"
)

// Unfollow ingrates
func (a *App) Unfollow() {
	log.Println("Collecting data...\n")
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

// Like and follow user's followers
func (a *App) ShadowUser(username string, hasNew bool) {
	if hasNew == true {
		user := a.GetUserByUsername(username)
		a.getUserFollowers(user)
		a.checkUserFollowers()
	}

	a.likeAndFollow()
}
