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
func (a *App) LikeFeed() {
	a.getFollowings()

}

// moreFollowers, _ := a.db2.ReadAll("followers")

// 	// iterate over morefish creating a new fish for each record
// 	followers := []InstagramUser{}
// 	for _, follower := range moreFollowers {
// 		f := InstagramUser{}
// 		json.Unmarshal([]byte(follower), &f)
// 		followers = append(followers, f)
// 	}

// 	fmt.Printf("It's a lot of fish! %#v\n", followers)
