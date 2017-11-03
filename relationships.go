package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// getFollowers gets users that follow us.
func (a *App) getFollowers() {
	log.Println("Collecting your 'Followers' list")

	// Resp data
	// Username: string
	// HasAnonymousProfilePicture: bool
	// ProfilePictureID: int
	// ProfilePictureURL:	URL
	// FullName: string
	// ID: int
	// IsVerified: bool
	// IsPrivate: bool
	// IsFavorite: bool
	// IsUnpublished: bool
	resp, err := a.api.SelfTotalUserFollowers()
	if err != nil {
		panic(err)
	}

	for _, user := range resp.Users {
		username := strings.ToLower(user.Username)

		log.Printf("USER data %+v \n", user)

		a.followers[username] = true
		a.db2.Write("followers", username, InstagramUser{ID: user.ID, Username: username, IsPrivate: user.IsPrivate, IsLiked: false})
	}
}

func (a *App) getUserFollowers(vip *InstagramUser) {
	fmt.Printf("Getting %s's followers 🐶 \n\n", vip.Username)

	resp, err := a.api.UserFollowers(vip.ID, "")
	if err != nil {
		panic(err)
	}

	// --- Structure ---
	// Username
	// HasAnonymousProfilePicture: false
	// ProfilePictureID: 1500417914300789823_408446133
	// ProfilePictureURL
	// FullName:
	// ID: number
	// IsVerified: false
	// IsPrivate: false
	// IsFavorite: false
	// IsUnpublished: false

	for _, user := range resp.Users {
		// log.Printf("USER data %+v \n", user)

		username := strings.ToLower(user.Username)
		currentUser := &InstagramUser{
			ID:        user.ID,
			Username:  username,
			IsPrivate: user.IsPrivate,
			IsLiked:   false,
		}

		a.db2.Write("user_followers", username, currentUser)
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
		a.db2.Write("followings", username, InstagramUser{ID: user.ID, Username: username, IsPrivate: user.IsPrivate, IsLiked: false})
	}
}

// check if user's followers are good to follow
func (a *App) checkUserFollowers() {
	collection := "user_followers"

	results, _ := a.db2.ReadAll(collection)

	data := []InstagramUser{}
	for _, user := range results {
		iu := InstagramUser{}
		json.Unmarshal([]byte(user), &iu)
		if iu.IsPrivate == false && iu.IsLiked == false {
			data = append(data, iu)
		}
	}

	fmt.Printf("There are %d followers to check\n", len(data))

	c := 0
	for _, follower := range data {
		c++
		fID := follower.ID
		fUsername := follower.Username

		resp, err := a.api.GetUserByID(fID)
		if err != nil {
			log.Panicf("Got error when checking Follower %s", fUsername)
		}

		// A very simple way to narrow down people that won't follow you back.
		// People w/ tons of followers and a few followings don't give a shit about you.
		isTooFamous := resp.User.FollowingCount < resp.User.FollowerCount
		log.Printf("[%s][%v/%v] %d followings, %d followers, famous? %t", fUsername, c, len(data), resp.User.FollowingCount, resp.User.FollowerCount, isTooFamous)

		// TODO: we set a bad follower as IsLiked: true, but there should be another property to describe that isn't good to follow
		uptUser := InstagramUser{ID: follower.ID, Username: follower.Username, IsPrivate: follower.IsPrivate, IsLiked: isTooFamous}
		if err := a.db2.Write(collection, follower.Username, uptUser); err != nil {
			fmt.Printf("Error while setting isLiked at true, %s", err)
		}

		time.Sleep(time.Duration(a.Wait) * time.Second)
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

	fmt.Printf("\n 🖕 🖕 🖕  Beginning Unfollow 🖕 🖕 🖕 \n\n")

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

func (a *App) follow() {
	// ID := ???
	// // TODO add a comment?
	// respFollow, errFollow := a.api.Follow(ID)
	// if errFollow != nil {
	// 	log.Panicf("Got error when Following : %s", errFollow)
	// }
	// log.Printf("Started to follow %s - response: %v ", f.Username, respFollow)
}
