package templates

import "architectSocial/app/repository"

type ListData struct {
	PageTitle     string
	Users         []repository.UserItem
	CurrentUserId string
	ShowSearch    bool
}
