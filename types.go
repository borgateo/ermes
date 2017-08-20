package main

import (
	"github.com/ahmdrz/goinsta"
	"gopkg.in/mgo.v2"
)

// Type App is the Follow-Sync Application.
type App struct {
	api        *goinsta.Instagram
	db         *mgo.Database
	session    *mgo.Session
	username   string
	password   string
	Wait       int             // How long to wait between instagram calls.
	followings map[string]bool // who are we following?
	followers  map[string]bool // who follows us?
	leeches    []string        // users we follow who don't follow us back
}

type InstagramUser struct {
	ID        int64
	Username  string
	IsPrivate bool
}

type Followers struct {
	ID int64 `bson:"id"`
	// Time      time.Time `bson:"time"`
	Username   string `bson:"username"`
	Following  int64  `bson:"following"`
	IsGood     bool   `bson:"isGood"`
	IsChecked  bool   `bson:"isChecked"`
	IsFollowed bool   `bson:"isFollowed"`
}

type Followings struct {
	ID       int64  `bson:"id"`
	Username string `bson:"username"`
	Me       string `bson:"me"`
}
