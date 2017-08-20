package main

import (
	"fmt"
	"os"

	"github.com/ahmdrz/goinsta"
	"gopkg.in/mgo.v2"
)

func (a *App) Login() {

	fmt.Printf("Login...\n")

	a.api = goinsta.New(a.username, a.password)

	if err := a.api.Login(); err != nil {
		fmt.Printf("Login error: %s\n", err)
	}
}

func (a *App) Logout() {
	fmt.Printf("Logout \n")
	a.api.Logout()
}

func (a *App) InitDB() {
	session, err := mgo.Dial(os.Getenv("MONGO_PORT"))
	if err != nil {
		panic(err)
	}
	a.session = session

	a.db = a.session.DB(os.Getenv("MONGO_DB"))
}

func (a *App) CloseDB() {
	a.session.Close()
}
