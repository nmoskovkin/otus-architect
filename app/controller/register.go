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

func CreateRegisterGetHandler(templ *template.Template) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := templ.ExecuteTemplate(w, "register.html", templates.RegisterData{
			PageTitle: "Register New User",
			Errors:    []string{},
		})
		if err != nil {
			return NewHTTPError(err, 500, "")
		}

		return nil
	}
}

func CreateRegisterPostHandler(templ *template.Template, db *sql.DB, sessionWrapper helpers.SessionWrapper) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		registerService := domain.CreateRegisterUserService(repository.CreateMysqlUserRepository(db))
		err := r.ParseForm()
		if err != nil {
			return NewHTTPError(err, 400, "")
		}
		dto := domain.RegisterUserDto{
			FirstName:            r.Form.Get("first_name"),
			LastName:             r.Form.Get("last_name"),
			Age:                  r.Form.Get("age"),
			Gender:               r.Form.Get("gender"),
			Interests:            r.Form.Get("interests"),
			City:                 r.Form.Get("city"),
			Password:             r.Form.Get("password"),
			PasswordConfirmation: r.Form.Get("password-confirmation"),
		}
		validationResult, userId, err := registerService(&dto)
		if err != nil {
			return NewHTTPError(err, 400, "")
		}

		if validationResult != nil {
			err := templ.ExecuteTemplate(w, "register.html", templates.RegisterData{
				PageTitle: "Register New User",
				Errors:    validationResult.GetAllErrors(),
				FirstName: dto.FirstName,
				LastName:  dto.LastName,
				Age:       dto.Age,
				Gender:    dto.Gender,
				Interests: dto.Interests,
				City:      dto.City,
			})
			if err != nil {
				return NewHTTPError(err, 500, "")
			}
			return nil
		}

		err = sessionWrapper.SetRegistrationId(userId, r, w)
		if err != nil {
			return NewHTTPError(err, 500, "")
		}
		http.Redirect(w, r, "/auth", 301)

		return nil
	}
}
