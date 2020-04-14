package controller

import (
	"architectSocial/app/templates"
	"architectSocial/domain"
	"html/template"
	"net/http"
)

func CreateRegisterGetHandler(templ *template.Template) *Handler {
	return NewHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			err := templ.ExecuteTemplate(w, "register.html", templates.RegisterData{
				PageTitle: "Register New User",
				Errors:    []string{},
			})
			if err != nil {
				return NewHTTPError(err, 500, "")
			}

			return nil
		},
		templ,
	)
}

func CreateRegisterPostHandler(registerUser domain.RegisterUserService, templ *template.Template) *Handler {
	return NewHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			err := r.ParseForm()
			if err != nil {
				return NewHTTPError(err, 500, "")
			}
			dto := domain.RegisterUserDto{
				FirstName: r.Form.Get("first_name"),
				LastName:  r.Form.Get("last_name"),
				Age:       r.Form.Get("age"),
				Gender:    r.Form.Get("gender"),
				Interests: r.Form.Get("interests"),
				City:      r.Form.Get("city"),
			}
			validationResult, err := registerUser(&dto)
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

			return nil
		},
		templ,
	)
}
