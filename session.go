package main

import (
	"fmt"
	"log"

	"github.com/ahmdrz/goinsta"
	"github.com/borteo/ermes/config"
	scribble "github.com/nanobox-io/golang-scribble"
)

func (a *App) Login() {
	fmt.Printf("🔑  Login - Please wait...\n\n")

	a.api = goinsta.New(a.username, a.password)

	if err := a.api.Login(); err != nil {
		log.Panicf("Login error: %s\n", err)
	}
}

func (a *App) Logout() {
	fmt.Printf("Logout \n")
	a.api.Logout()
}

func (a *App) InitDB() {
	db, err := scribble.New(config.DATA_PATH, nil)
	if err != nil {
		log.Panicf("Error initializing DB: %s", err)
	}

	a.db2 = db
}
