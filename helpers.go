package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// sortKeys returns the sorted keys from the follower/ing lists.
func (a *App) sortKeys(dict map[string]bool) []string {
	var keys []string
	for key, _ := range dict {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// max returns the greater of two ints (math.Max does float64)
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getUserId gets the Instagram user PK ID (an int64) as a string.
func (a *App) getUserId(username string) string {
	user, err := a.api.GetUserByUsername(username)
	if err != nil {
		log.Panicf("Can't getUserId %s: %s", username, err)
	}

	return strconv.Itoa(int(user.User.ID))
}

// compareLists compares the following to the followers.
func (a *App) compareLists() {
	// See who I am following that doesn't love me back.
	for username, _ := range a.followings {
		if _, ok := a.followers[username]; !ok {
			a.leeches = append(a.leeches, username)
		}
	}
}

func random(val int) int {
	min := val - val/2
	max := val + val/2
	rand.Seed(time.Now().Unix())

	return rand.Intn(max-min) + min
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func askForConfirmation(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}
