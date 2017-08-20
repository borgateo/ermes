package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (a *App) FollowVIPFollowers(username string) {

	resp, err := a.api.GetUserByUsername(username)
	if err != nil {
		panic(err)
	}

	vip := &InstagramUser{
		ID:        resp.User.ID,
		Username:  resp.User.Username,
		IsPrivate: resp.User.IsPrivate,
	}
	fmt.Printf("VIP: %v \n", vip)
	fmt.Printf("GetUserByUsername: %v \n", resp)

	a.GetVIPFollowers(vip)
	time.Sleep(time.Duration(a.Wait) * time.Second)
	a.checkFollowers(vip)
	a.likeAndFollowFollowers()
}

// Save all user's followers in the DB
// in the "followers" collection
func (a *App) GetVIPFollowers(vip *InstagramUser) {
	fmt.Printf("\n--- Fetch %s's followers 🐶 ---\n\n", vip.Username)

	resp, err := a.api.UserFollowers(vip.ID, "")
	if err != nil {
		panic(err)
	}

	// --- Structure ---
	// Username
	// HasAnonymousProfilePicture: false
	// ProfilePictureID:1500417914300789823_408446133
	// ProfilePictureURL
	// FullName:
	// ID:408446133
	// IsVerified:false
	// IsPrivate:false
	// IsFavorite:false
	// IsUnpublished:false

	c := a.db.C("followers")

	for _, user := range resp.Users {
		count, err := c.Find(bson.M{"id": user.ID}).Count()
		if err != nil {
			log.Println("Error Finding Profile: ", err.Error())
		}

		if count > 0 {
			log.Printf("User %v already exists", user.ID)
			continue
		}

		// set `isChecked` to `isPrivate` cos
		// if it's private we can't see the user medias
		// we should just skip it
		err = c.Insert(
			&Followers{
				ID:         user.ID,
				Username:   user.Username,
				Following:  vip.ID,
				IsGood:     false,
				IsChecked:  user.IsPrivate,
				IsFollowed: false,
			})
		if err != nil {
			log.Printf("Error creating Profile: %s", err.Error())
		}

		fmt.Printf("\n- adding: %s", user.Username)

	}
}

func (a *App) checkFollowers(vip *InstagramUser) {
	fmt.Printf("\n--- Check if %s's followers are good to follow ---\n", vip.Username)
	c := a.db.C("followers")

	var results []Followers
	err := c.Find(bson.M{"isChecked": false, "isFollowed": false}).All(&results)

	if err != nil {
		panic(err)
	}

	for follower := range results {
		fID := results[follower].ID
		fUsername := results[follower].Username

		resp, err := a.api.GetUserByID(fID)
		if err != nil {
			log.Panicf("Got error when checking Follower %s", fUsername)
		}
		isGood := resp.User.FollowingCount > resp.User.FollowerCount
		fmt.Printf("\n %s: followings %d, followers %d 👍 %t", fUsername, resp.User.FollowingCount, resp.User.FollowerCount, isGood)
		err = c.Update(bson.M{"id": fID}, bson.M{"$set": bson.M{"isGood": isGood, "isChecked": true}})

		time.Sleep(time.Duration(a.Wait) * time.Second)
	}
}

// feed of good followers
func (a *App) likeAndFollowFollowers() {
	fmt.Printf("\n--- 👍 Like good followers' media and follow them ❤️ ---\n")
	c := a.db.C("followers")

	var results []Followers
	err := c.Find(bson.M{"isGood": true, "isFollowed": false}).All(&results)

	if err != nil {
		panic(err)
	}

	for follower := range results {
		fID := results[follower].ID
		fUsername := results[follower].Username

		resp, err2 := a.api.UserFeed(fID, "", "")
		if err != nil {
			log.Panicf("Got error when getting UserFeed: %s", err2)
		}

		log.Printf("Followers feed: %+v", resp)

		// UserFeed response struct:
		// {Status:ok ,NumResults:13, AutoLoadMoreEnabled:true,
		// Items: [{TakenAt, ID, HasLiked, ...more }]
		counter := 0
		for _, item := range resp.Items {
			// next pic is already liked
			if item.HasLiked == true {
				continue
			}

			respLike, errLike := a.api.Like(item.ID)
			if errLike != nil {
				log.Panicf("Got error when Liking : %s", errLike)
			}
			log.Printf("Liked %+v", respLike)

			// Don't like more than 3 pics -- TODO: configurable?
			if counter > 3 {
				// and finally follow her/him
				// TODO add a comment?
				respFollow, errFollow := a.api.Follow(fID)
				if errFollow != nil {
					log.Panicf("Got error when Following : %s", errFollow)
				}
				log.Printf("Started to follow %s - response: %v ", fUsername, respFollow)
				break
			}

			counter++
			time.Sleep(time.Duration(a.Wait) * time.Second)
		}

		time.Sleep(time.Duration(a.Wait) * time.Second)
	}
}
