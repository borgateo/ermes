package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func (a *App) LikeFeed() {
	a.fetchFollowers()
}

func (a *App) fetchFollowers() {
	log.Println("Fetch 'Followers'")
	resp, err := a.api.SelfTotalUserFollowers()
	if err != nil {
		panic(err)
	}

	c := a.db.C("followings")

	for _, user := range resp.Users {
		count, err := c.Find(bson.M{"id": user.ID}).Count()
		if err != nil {
			log.Println("Error Finding Profile: ", err.Error())
		}

		if count > 0 {
			log.Printf("User %v already exists", user.ID)
			continue
		}

		err = c.Insert(
			&Followings{
				ID:       user.ID,
				Username: user.Username,
				Me:       a.username,
			})
		if err != nil {
			log.Printf("Error creating Profile: %s", err.Error())
		}

		fmt.Printf("\n- adding: %s", user.Username)
	}

}

func (a *App) likeFollowings() {
	fmt.Printf("\n--- 👍 Like followings media 👍 ---\n")
	c := a.db.C("followings")

	var results []Followings
	err := c.Find(bson.M{"Me": a.username}).All(&results)

	if err != nil {
		panic(err)
	}

	for following := range results {
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
