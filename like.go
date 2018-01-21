package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ahmdrz/goinsta/response"
	"github.com/borteo/ermes/config"
	humanize "github.com/dustin/go-humanize"
)

func (a *App) likeFeed(types string, isLimited bool) {
	a.likeAndFollow(types, false, isLimited)
}

func (a *App) likeAndFollowFeed(types string, isLimited bool) {
	a.likeAndFollow(types, true, isLimited)
}

func (a *App) getFilteredData(types string, isLimited bool) []InstagramUser {
	limit := config.FOLLOWERS_MAX
	all, _ := a.db2.ReadAll(types)

	// filter by non private, not liked yet
	data := []InstagramUser{}
	for _, user := range all {
		iu := InstagramUser{}
		json.Unmarshal([]byte(user), &iu)
		if iu.IsPrivate == false && iu.IsLiked == false && iu.IsGood == true {
			data = append(data, iu)
		}
	}

	if isLimited && len(data) > limit {
		data = data[0:limit]
	}

	fmt.Printf("\nYou have %d saved, and filtered %d %v\n", len(all), len(data), types)
	return data
}

func (a *App) likeAndFollow(types string, shouldFollow bool, isLimited bool) {
	var (
		nUsers        = 0
		nMedia        = 0
		feedErrors    = 0
		feedErrorsMax = 10
		likeMax       = config.LIKE_MAX
	)

	data := a.getFilteredData(types, isLimited)
	fmt.Printf("%d users to process\n\n", len(data))

	for _, user := range data {
		nUsers++
		fmt.Printf("\n")
		log.Printf("\n")
		fmt.Printf("Progress: %d/%d (%.2f%%) \n", nUsers, len(data), float64(nUsers)/float64(len(data))*float64(100))
		var delaySecs time.Duration = time.Duration(config.WAITING_TIME*(len(data)-nUsers)) * time.Second
		fmt.Printf("⏱  %s \n", humanize.Time(time.Now().Add(delaySecs)))

		fmt.Printf("💕  Spreading love to '%+v'\n", user.Username)

		resp, err := a.api.UserFeed(user.ID, "", "")
		if err != nil {
			log.Printf("ERROR: on 'UserFeed' %s", err)
			if feedErrors == feedErrorsMax {
				log.Panicf("PANIC: got too many errors when fetching userFeed")
			} else {
				feedErrors++
				continue
			}
		}
		// reset feed errors count if reaches this point
		feedErrors = 0

		// log.Printf("\nFeed: %+v", resp)

		// UserFeed response struct:
		// Status:ok, NumResults: int, AutoLoadMoreEnabled: bool,
		// Items: [{TakenAt, ID, HasLiked, ...more }]
		nMedia = 0
		for _, item := range resp.Items {
			// Don't like more than N pics
			if nMedia == likeMax || nMedia == len(resp.Items) {
				// When I reach the limit, I can follow the user
				if shouldFollow == true {
					_, errFollow := a.api.Follow(user.ID)
					if errFollow != nil {
						log.Printf("Got error on Follow: %s", errFollow)
						continue
					}
					log.Printf("Started to follow %s", user.Username)

					uptUser := InstagramUser{
						ID:        user.ID,
						Username:  user.Username,
						IsPrivate: user.IsPrivate,
						IsGood:    user.IsGood,
						IsChecked: user.IsChecked,
						IsLiked:   true,
					}
					if err := a.db2.Write(types, user.Username, uptUser); err != nil {
						fmt.Printf("Error while setting isLiked at true, %s", err)
					}
				}
				break
			}
			// Move on if pic already liked
			if item.HasLiked == true {
				continue
			}

			_, errLike := a.api.Like(item.ID)
			if errLike != nil {
				log.Printf("Got error when Liking : %s", errLike)
				continue
			}

			nMedia++
			log.Printf("👍  %v's media: [%v]", user.Username, item.ID)
			time.Sleep(time.Duration(random(config.WAITING_TIME)) * time.Second)
		}

		// set 'isLiked' at true
		uptUser := InstagramUser{
			ID:        user.ID,
			Username:  user.Username,
			IsPrivate: user.IsPrivate,
			IsGood:    true,
			IsLiked:   true,
		}
		if err := a.db2.Write(types, user.Username, uptUser); err != nil {
			fmt.Printf("Error while setting isLiked at true, %s", err)
		}

		time.Sleep(time.Duration(random(config.WAITING_TIME)) * time.Second)
	}
}

func (a *App) likeTimeline(timeline response.FeedsResponse) {
	for _, item := range timeline.Items {
		if item.HasLiked == true {
			log.Printf("%v's media already liked: [%v]", item.User.Username, item.ID)
			continue
		}

		_, errLike := a.api.Like(item.ID)
		if errLike != nil {
			log.Printf("Got error when Liking : %s", errLike)
			continue
		}

		log.Printf("👍  %v's media: [%v]", item.User.Username, item.ID)
		time.Sleep(time.Duration(random(config.WAITING_TIME)) * time.Second)
	}

}
