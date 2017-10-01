package main

import (
	"log"
	"time"
)

func (a *App) likeFollowings() {
	for following := range a.followings {
		fID := results[following].ID
		fUsername := results[following].Username

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
