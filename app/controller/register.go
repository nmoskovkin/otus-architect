package controller

import (
	"architectSocial/app/repository"
	"architectSocial/app/templates"
	"architectSocial/domain"
	"database/sql"
	"html/template"
	"net/http"
)

func CreateRegisterHandler(templ *template.Template, db *sql.DB) *Handler {
	return NewHandler(
		func(writer http.ResponseWriter, request *http.Request) error {
			var handler ErrorReturningHandlerFunc
			if request.Method == http.MethodPost {
				handler = createRegisterPostHandler(domain.CreateRegisterUserService(repository.CreateUserRepository(db)), templ)
			} else {
				handler = createRegisterGetHandler(templ)
			}

			return handler(writer, request)
		},
		templ,
	)

	//func(writer http.ResponseWriter, request *http.Request) error {
	//	var handler func(writer http.ResponseWriter, request *http.Request)
	//	if request.Method == http.MethodPost {
	//		handler = createRegisterPostHandler(domain.CreateRegisterUserService(repository.CreateUserRepository(db)), templ)
	//	} else {
	//		handler = createRegisterGetHandler(templ)
	//	}
	//
	//	handler(writer, request)
	//}
	//return func(writer http.ResponseWriter, request *http.Request) {
	//	var handler func(writer http.ResponseWriter, request *http.Request)
	//	if request.Method == http.MethodPost {
	//		handler = createRegisterPostHandler(domain.CreateRegisterUserService(repository.CreateUserRepository(db)), templ)
	//	} else {
	//		handler = createRegisterGetHandler(templ)
	//	}
	//
	//	handler(writer, request)
	//}
}

//func CreateRegisterHandler(registerUser domain.RegisterUserService, templ *template.Template) Handler {
//
//}

func createRegisterGetHandler(templ *template.Template) ErrorReturningHandlerFunc {
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

func createRegisterPostHandler(registerUser domain.RegisterUserService, templ *template.Template) ErrorReturningHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
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
	}
}
