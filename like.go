package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func (a *App) likeFollowersPosts() {
	a.likePosts("followers")
}

func (a *App) likeFollowingsPosts() {
	a.likePosts("followings")
}

func (a *App) likePosts(types string) {
	likeMaxString := os.Getenv("LIKE_MAX")
	likeMax, err := strconv.Atoi(likeMaxString)
	if err != nil {
		likeMax = 3
	}

	// get all users
	all, _ := a.db2.ReadAll(types)

	// filter by non private and not liked yet
	data := []InstagramUser{}
	for _, user := range all {
		iu := InstagramUser{}
		json.Unmarshal([]byte(user), &iu)
		if iu.IsPrivate == false && iu.IsLiked == false {
			data = append(data, iu)
		}
	}

	fmt.Printf("You have %d %v; To like %d users \n", len(all), types, len(data))

	nUsers := 0
	feedErrors := 0
	for _, following := range data {
		nUsers++
		fmt.Printf("\nProgress %d/%d (%f%%) \n", nUsers, len(data), float64(nUsers)/float64(len(data))*float64(100))

		fmt.Printf("\n")
		log.Printf("💕  Spreading some love to '%+v'", following.Username)

		resp, err := a.api.UserFeed(following.ID, "", "")
		if err != nil {
			log.Printf("ERROR: on 'UserFeed' %s", err)
			if feedErrors >= 5 {
				log.Panicf("PANIC: got too many errors when fetching userFeed")
			} else {
				feedErrors++
				continue
			}
		}

		// reset feed errors count if reaches this point
		feedErrors = 0

		// log.Printf("\nFeed: %+v", resp)

		// set 'isLiked' at true
		uptUser := InstagramUser{ID: following.ID, Username: following.Username, IsPrivate: following.IsPrivate, IsLiked: true}
		if err := a.db2.Write(types, following.Username, uptUser); err != nil {
			fmt.Printf("Error while setting isLiked at true, %s", err)
		}

		// UserFeed response struct:
		// Status:ok, NumResults: int, AutoLoadMoreEnabled: bool,
		// Items: [{TakenAt, ID, HasLiked, ...more }]
		n := 0
		for _, item := range resp.Items {
			// Don't like more than N pics
			if n >= likeMax {
				break
			}
			// Move on if pic already liked
			if item.HasLiked == true {
				continue
			}

			respLike, errLike := a.api.Like(item.ID)
			if errLike != nil {
				log.Printf("ERROR: on 'like' %s", errLike)
				continue
			}

			n++
			log.Printf("👍  %+v", respLike)
			time.Sleep(time.Duration(a.Wait) * time.Second)
		}

		time.Sleep(time.Duration(a.Wait) * time.Second)
	}
}

// feed of good followers
func (a *App) likeAndFollow() {
	fmt.Printf("\nLike good followers' media and follow\n")

	collection := "user_followers"
	likeMaxString := os.Getenv("LIKE_MAX")
	likeMax, err := strconv.Atoi(likeMaxString)
	if err != nil {
		likeMax = 3
	}

	all, _ := a.db2.ReadAll(collection)

	// filter by non private and not liked yet
	data := []InstagramUser{}
	for _, user := range all {
		iu := InstagramUser{}
		json.Unmarshal([]byte(user), &iu)
		if iu.IsPrivate == false && iu.IsLiked == false {
			data = append(data, iu)
		}
	}

	fmt.Printf("There are %d users to follow and like\n", len(data))

	for _, follower := range data {
		fID := follower.ID
		fUsername := follower.Username

		resp, err := a.api.UserFeed(fID, "", "")
		if err != nil {
			fmt.Printf("Error while getting UserFeed, %s", err)
			continue
		}

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

			if counter >= likeMax {
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
