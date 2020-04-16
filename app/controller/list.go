package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"database/sql"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

func CreateListGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		isAuth, id, err := sessionWrapper.IsAuthenticated(r)
		var curentUserId string
		if isAuth && err == nil {
			curentUserId = id
		}

		userRepository := repository.CreateMysqlUserRepository(db)
		userList, err := userRepository.GetAll(repository.GetAllFilter{})
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		listData := templates.ListData{
			Users:         userList,
			PageTitle:     "All People",
			CurrentUserId: curentUserId,
		}
		err = templ.ExecuteTemplate(w, "list.html", listData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}

func CreateFriendsListGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		isAuth, id, err := sessionWrapper.IsAuthenticated(r)
		var currentUserId string
		if isAuth && err == nil {
			currentUserId = id
		} else {
			return NewHTTPError(err, 403, "")
		}

		userRepository := repository.CreateMysqlUserRepository(db)
		friendsRepository := repository.CreateMysqlFriendsRepository(db)

		currentUserUUID, err := uuid.Parse(currentUserId)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		list, err := friendsRepository.GetFriends(currentUserUUID)
		userList, err := userRepository.GetAll(repository.GetAllFilter{Ids: list})
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		listData := templates.ListData{
			Users:         userList,
			PageTitle:     "Friends",
			CurrentUserId: currentUserId,
		}
		err = templ.ExecuteTemplate(w, "list.html", listData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}
