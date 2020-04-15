package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

func CreateAuthGetHandler(templ *template.Template, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, err := sessionWrapper.GetRegistrationId(r)
		templData := templates.AuthData{
			PageTitle: "Register New User",
		}
		if err == nil {
			templData.Id = userId
		}
		err = templ.ExecuteTemplate(w, "auth.html", templData)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}

func CreateAuthPostHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		repo := repository.CreateMysqlUserRepository(db)
		err := r.ParseForm()
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		id := r.Form.Get("id")
		userPassword := r.Form.Get("password")
		uuid_, err := uuid.Parse(id)
		if err != nil {
			// error incorrect login or email
			return NewHTTPError(err, 500, "")
		}
		users, err := repo.FindById(uuid_)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		if len(users) == 0 {
			// error user not found
		}

		//TODO move to service
		password := users[0]["password"]
		salt := users[0]["salt"]

		passwordValue, okV := password.([]byte)
		saltValue, okS := salt.([]byte)
		if !okS || !okV {
			return NewHTTPError(errors.New("!okS || !okV "), 500, "")
		}
		err = bcrypt.CompareHashAndPassword([]byte(passwordValue), append([]byte(userPassword), saltValue...))
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		fmt.Println(users)

		return nil
	}
}
