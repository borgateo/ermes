package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)

func (a *App) likeFollowersPosts() {
	likeMaxString := os.Getenv("LIKE_MAX")
	likeMax, err := strconv.Atoi(likeMaxString)
	if err != nil {
		likeMax = 3
	}

	data, _ := a.db2.ReadAll("followers")

	// followings := []InstagramUser{}
	for _, following := range data {
		f := InstagramUser{}
		json.Unmarshal([]byte(following), &f)

		log.Printf("\nSending some love to %+v", f.Username)

		resp, err := a.api.UserFeed(f.ID, "", "")
		if err != nil {
			log.Panicf("Got error when getting UserFeed: %s", err)
		}

		// log.Printf("\nFeed: %+v", resp)

		// UserFeed response struct:
		// {Status:ok ,NumResults:13, AutoLoadMoreEnabled:true,
		// Items: [{TakenAt, ID, HasLiked, ...more }]
		n := 0
		for _, item := range resp.Items {
			// Don't like more than N pics
			if n > likeMax {
				break
			}
			// Move on if pic already liked
			if item.HasLiked == true {
				continue
			}

			respLike, errLike := a.api.Like(item.ID)
			if errLike != nil {
				log.Panicf("Got error when Liking : %s", errLike)
			}
			log.Printf("👍 liked %+v", respLike)

			n++
			time.Sleep(time.Duration(a.Wait) * time.Second)
		}

		time.Sleep(time.Duration(a.Wait) * time.Second)
	}
}
