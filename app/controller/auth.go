package controller

import (
	"architectSocial/app/helpers"
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"architectSocial/domain"
	"database/sql"
	"html/template"
	"net/http"
)

func CreateAuthGetHandler(templ *template.Template, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, err := sessionWrapper.GetRegistrationId(r)
		templData := templates.AuthData{
			PageTitle: "Authenticate",
		}
		if err == nil {
			templData.Login = userId
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
		authService := domain.CreateAuthUserService(repository.CreateMysqlUserRepository(db))
		err := r.ParseForm()
		if err != nil {
			return NewHTTPError(err, 400, "")
		}
		authUserDto := &domain.AuthUserDto{
			Login:    r.Form.Get("login"),
			Password: r.Form.Get("password"),
		}
		validationResult, isAuthenticated, err := authService(authUserDto)
		if err != nil {
			return NewHTTPError(err, 400, "")
		}
		if validationResult != nil {
			err := templ.ExecuteTemplate(w, "auth.html", templates.AuthData{
				PageTitle: "Authenticate",
				Errors:    validationResult.GetAllErrors(),
				Login:     r.Form.Get("login"),
			})
			if err != nil {
				return NewHTTPError(err, 500, "")
			}
			return nil
		}
		if isAuthenticated {
			err := sessionWrapper.SetAuthenticated(authUserDto.Login, r, w)
			if err != nil {
				return NewHTTPError(err, 500, "")
			}

			return nil
		}

		err = templ.ExecuteTemplate(w, "auth.html", templates.AuthData{
			PageTitle: "Authenticate",
			Errors:    []string{"invalid credentials"},
			Login:     r.Form.Get("login"),
		})
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		http.Redirect(w, r, "/list", 301)
		return nil
	}
}
