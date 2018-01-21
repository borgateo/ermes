package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ahmdrz/goinsta/response"
	"github.com/borteo/ermes/config"
	humanize "github.com/dustin/go-humanize"
)

// getFollowing gets users that we follow.
func (a *App) getFollowings() error {
	fmt.Println("Collecting your 'Followings' list...")
	resp, err := a.api.SelfTotalUserFollowing()
	if err != nil {
		return err
	}

	for _, user := range resp.Users {
		username := strings.ToLower(user.Username)
		a.followings[username] = true
		a.db2.Write("followings", username, InstagramUser{
			ID:        user.ID,
			Username:  username,
			IsPrivate: user.IsPrivate,
			IsLiked:   false,
			IsChecked: true,
			IsGood:    !user.IsPrivate,
		})
	}
	return nil
}

// getFollowers gets users that follow us.
func (a *App) getFollowers() error {
	log.Println("Collecting your 'Followers' list...")

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
		return err
	}

	for _, user := range resp.Users {
		username := strings.ToLower(user.Username)

		uptUser := InstagramUser{
			ID:        user.ID,
			Username:  username,
			IsPrivate: user.IsPrivate,
			IsLiked:   false,
			IsGood:    true,
			IsChecked: true,
		}
		a.followers[username] = true
		a.db2.Write("followers", username, uptUser)
	}
	return nil
}

// Get passed user.ID followers.
// @param isLimited: if true fetch only the amount of pages defined
func (a *App) fetchUserFollowers(vipID int64, isLimited bool) (response.UsersResponse, error) {
	usersLimit := 1500
	resp := response.UsersResponse{}

	for {
		tempResp, err := a.api.UserFollowing(vipID, resp.NextMaxID)
		if err != nil {
			return response.UsersResponse{}, err
		}
		resp.Users = append(resp.Users, tempResp.Users...)
		resp.PageSize += tempResp.PageSize
		if !tempResp.BigList || (isLimited && len(resp.Users) >= usersLimit) {
			return resp, nil
		}
		resp.NextMaxID = tempResp.NextMaxID
		resp.Status = tempResp.Status
	}
}

func (a *App) getUserFollowers(vip *InstagramUser, isLimited bool) error {
	fmt.Printf("Collecting %s's followers 🐶 \n\n", vip.Username)

	collection := "user_followers_" + vip.Username
	resp, err := a.api.TotalUserFollowing(vip.ID)
	// Use it to limit the number of pages requested:
	// resp, err := a.fetchUserFollowers(vip.ID, isLimited)
	if err != nil {
		return err
	}

	// --- resp structure ---
	// BigList: false
	// PageSize: 200
	// Users: []
	// log.Printf("There are %v \n", len(resp.Users))
	// log.Printf("There are %+v \n", resp)

	for _, user := range resp.Users {
		// --- user structure ---
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
		// log.Printf("USER data %+v \n", user)

		username := strings.ToLower(user.Username)
		currentUser := &InstagramUser{
			ID:        user.ID,
			Username:  username,
			IsPrivate: user.IsPrivate,
			IsLiked:   false,
			IsChecked: false,
			IsGood:    !user.IsPrivate,
		}

		a.db2.Write(collection, username, currentUser)
	}
	return nil
}

// check if user's followers are good to follow
func (a *App) checkUserFollowers(username string) {
	limit := 2000
	collection := "user_followers_" + username

	results, _ := a.db2.ReadAll(collection)

	data := []InstagramUser{}
	for _, user := range results {
		iu := InstagramUser{}
		json.Unmarshal([]byte(user), &iu)
		if iu.IsPrivate == false && iu.IsChecked == false {
			data = append(data, iu)
		}
	}
	totalUsersLen := len(data)

	// Don't analyze more than 'limit' amount of users
	// Why? it requires 3 hours to analyze 1k users.
	// Ideally we could analyze 1k, then process the outcome and restart the loop
	if len(data) > limit {
		data = data[0:limit]
	}

	fmt.Printf("🔍  There are %d followers to check; %d more available \n", len(data), totalUsersLen-len(data))
	var delaySecs time.Duration = time.Duration(config.WAITING_TIME*len(data)) * time.Second
	fmt.Printf("⏱  %s \n\n", humanize.Time(time.Now().Add(delaySecs)))

	counter := 0
	for _, follower := range data {
		counter++
		fID := follower.ID
		fUsername := follower.Username

		resp, err := a.api.GetUserByID(fID)
		if err != nil {
			log.Printf("Got error checking %s. Moving on...", fUsername)
			continue
		}

		// A very simple way to narrow down people that won't follow you back.
		// People w/ tons of followers and a few followings don't give a shit about you.
		isGood := resp.User.FollowingCount > resp.User.FollowerCount
		log.Printf("[%v/%v][%s] %d followings, %d followers, is good? %t", counter, len(data), fUsername, resp.User.FollowingCount, resp.User.FollowerCount, isGood)

		uptUser := InstagramUser{
			ID:        follower.ID,
			Username:  follower.Username,
			IsPrivate: follower.IsPrivate,
			IsLiked:   false,
			IsGood:    isGood,
			IsChecked: true,
		}
		if err := a.db2.Write(collection, follower.Username, uptUser); err != nil {
			fmt.Printf("Error while setting isLiked at true, %s", err)
		}

		time.Sleep(time.Duration(config.WAITING_TIME) * time.Second)
	}
}

// unfollow leeches -- TODO pass array to unfollow
func (a *App) unfollowLeeches() {
	var (
		counter   = 1
		remaining = len(a.leeches)
	)

	if remaining == 0 {
		fmt.Printf("Congrats! You have no ingrates\n")
		return
	}

	fmt.Printf("\n 🖕 🖕 🖕 unfollow 🖕 🖕 🖕 \n\n")

	for _, username := range a.leeches {
		if _, ok := a.followings[username]; !ok {
			fmt.Printf("[ERROR] Username %s not found in following map\n", username)
			continue
		}

		// Unfollow.
		userIDStr := a.getUserId(username)
		randomInt := random(config.WAITING_TIME)
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
