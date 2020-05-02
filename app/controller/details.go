package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"architectSocial/domain"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

func CreateDetailsGetHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		isAuth, currentUserId, err := sessionWrapper.IsAuthenticated(r)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		if !isAuth {
			return NewHTTPError(fmt.Errorf("forbidden"), 403, "")
		}
		pageUserId := r.URL.Query().Get("id")
		if pageUserId == "" {
			return NewHTTPError(fmt.Errorf("user not found"), 404, "")
		}

		userRepository := repository.CreateMysqlUserRepository(db)
		friendsRepository := repository.CreateMysqlFriendsRepository(db)
		userList, err := userRepository.GetAll(repository.GetAllFilter{Id: pageUserId}, 1, 1)
		if len(userList) == 0 {
			return NewHTTPError(fmt.Errorf("user not found"), 404, "")
		}

		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		currentUserUUID, errCurrent := uuid.Parse(currentUserId)
		pageUserUUID, errPage := uuid.Parse(pageUserId)
		if errCurrent != nil || errPage != nil {
			return NewHTTPError(err, 500, "")
		}
		areFriends, err := friendsRepository.AreFriends(currentUserUUID, pageUserUUID)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		detailsData := templates.DetailsData{
			User:            userList[0],
			PageTitle:       "Users",
			IsAlreadyFriend: areFriends,
			CurrentUserId:   currentUserId,
		}
		err = templ.ExecuteTemplate(w, "details.html", detailsData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}

func CreateDetailsPostHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		pageUserId := r.URL.Query().Get("id")
		if pageUserId == "" {
			return NewHTTPError(fmt.Errorf("user not found"), 404, "")
		}
		userRepository := repository.CreateMysqlUserRepository(db)
		// TODO replace with a call of correct function, e.g. Get()
		userList, err := userRepository.GetAll(repository.GetAllFilter{Id: pageUserId}, 1, 1)
		if len(userList) == 0 {
			return NewHTTPError(fmt.Errorf("user not found"), 404, "")
		}

		isAuth, currentUserId, err := sessionWrapper.IsAuthenticated(r)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		if !isAuth {
			return NewHTTPError(fmt.Errorf("forbidden"), 403, "")
		}

		friendsRepository := repository.CreateMysqlFriendsRepository(db)
		makeFriendsService := domain.CreateMakeFriendService(friendsRepository)
		err = r.ParseForm()
		if err != nil {
			return NewHTTPError(err, 400, "")
		}
		makeFriendDto := &domain.MakeFriendDto{
			From: currentUserId,
			To:   pageUserId,
		}
		validationResult, err := makeFriendsService(makeFriendDto)
		if err != nil {
			return NewHTTPError(err, 400, "")
		}
		currentUserUUID, errCurrent := uuid.Parse(currentUserId)
		pageUserUUID, errPage := uuid.Parse(pageUserId)
		if errCurrent != nil || errPage != nil {
			return NewHTTPError(err, 500, "")
		}

		areFriends, err := friendsRepository.AreFriends(currentUserUUID, pageUserUUID)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		if validationResult != nil {
			detailsData := templates.DetailsData{
				User:            userList[0],
				PageTitle:       "Users",
				Errors:          validationResult.GetAllErrors(),
				IsAlreadyFriend: areFriends,
				CurrentUserId:   currentUserId,
			}
			err = templ.ExecuteTemplate(w, "details.html", detailsData)
			if err != nil {
				return NewHTTPError(err, 500, "")
			}
			return nil
		}

		http.Redirect(w, r, "/details?id="+pageUserId, 302)
		return nil
	}
}
