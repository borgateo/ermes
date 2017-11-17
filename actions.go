package main

import (
	"log"
	"os"

	"github.com/borteo/ermes/config"
)

// Unfollow ingrates
func (a *App) Unfollow() {
	log.Println("Collecting data\n")
	a.getFollowers()
	a.getFollowings()

	// Compare data.
	log.Println("Comparing followers and followings ✅ \n")
	a.compareLists()
	a.showList()

	a.unfollowLeeches()
}

// Like followers's feed
func (a *App) LikeFeedFollowers(skip bool) {
	if skip != true {
		// request user interaction
		c := askForConfirmation("All your stored Followers will be removed. Do you really want to continue?", 3)
		if c == false {
			return
		}
		os.RemoveAll(config.DATA_PATH + config.FOLLOWERS)
		a.getFollowers()
	}
	a.likeFeed(config.FOLLOWERS)
}

// Like followings's feed
func (a *App) LikeFeedFollowings(skip bool) {
	if skip != true {
		// request user interaction
		c := askForConfirmation("All your stored Followings will be removed. Do you really want to continue?", 3)
		if c == false {
			return
		}
		os.RemoveAll(config.DATA_PATH + config.FOLLOWINGS)
		a.getFollowings()
	}
	a.likeFeed(config.FOLLOWINGS)
}

// Like and follow user's followers
func (a *App) ShadowUser(username string, skip bool) {
	if skip != true {
		user := a.GetUserByUsername(username)
		a.getUserFollowers(user)
	}

	a.checkUserFollowers(username)
	a.likeAndFollowFeed(config.USER_FOLLOWERS + username)
}

// Like my feed
func (a *App) LikeMyTimeline() {
	page1, err := a.api.Timeline("")
	if err != nil {
		panic(err)
	}
	a.likeTimeline(page1)

	// I just want to like 2 pages
	if page1.MoreAvailable {
		page2, err := a.api.Timeline(page1.NextMaxID)
		if err != nil {
			panic(err)
		}
		a.likeTimeline(page2)
	}
}
