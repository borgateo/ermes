package main

import (
	"fmt"
	"os"

	"github.com/ahmdrz/goinsta"
	scribble "github.com/nanobox-io/golang-scribble"
)

func (a *App) Login() {

	fmt.Printf("🔑  Login - Please wait...\n")

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
	db, err := scribble.New(os.Getenv("STORAGE_PATH"), nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	a.db2 = db
}
