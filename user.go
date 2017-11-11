package main

func (a *App) GetUserByUsername(username string) *InstagramUser {

	resp, err := a.api.GetUserByUsername(username)
	if err != nil {
		panic(err)
	}

	user := &InstagramUser{
		ID:        resp.User.ID,
		Username:  resp.User.Username,
		IsPrivate: resp.User.IsPrivate,
		IsChecked: false,
		IsLiked:   false,
		IsGood:    false,
	}

	return user
}
