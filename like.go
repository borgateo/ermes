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

	data, _ := a.db2.ReadAll(types)

	fmt.Printf("You have %d %v\n", len(data), types)

	nUsers := 0
	for _, following := range data {
		nUsers++
		fmt.Printf("\nProgress %d/%d (%f%%) \n", nUsers, len(data), float64(nUsers)/float64(len(data))*float64(100))

		f := InstagramUser{}
		json.Unmarshal([]byte(following), &f)

		fmt.Printf("\n")
		log.Printf("💕  Spreading some love to '%+v'", f.Username)

		// skip private users
		if f.IsPrivate == true {
			log.Printf("WARNING: private user! Moving on...")
			continue
		}

		// skip already liked users
		if f.IsLiked == true {
			log.Printf("Already liked! Moving on...")
			continue
		}

		resp, err := a.api.UserFeed(f.ID, "", "")
		if err != nil {
			log.Panicf("Got error when getting UserFeed: %s", err)
		}

		// log.Printf("\nFeed: %+v", resp)

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
				log.Panicf("Got error when 'Liking': %s", errLike)
			}

			n++
			log.Printf("👍  %+v", respLike)
			time.Sleep(time.Duration(a.Wait) * time.Second)
		}

		// set 'isLiked' at true
		uptUser := InstagramUser{ID: f.ID, Username: f.Username, IsPrivate: f.IsPrivate, IsLiked: true}
		if err := a.db2.Write(types, f.Username, uptUser); err != nil {
			fmt.Println("Error", err)
		}

		time.Sleep(time.Duration(a.Wait) * time.Second)
	}
}
