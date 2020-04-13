package main

import (
	"architectSocial/model"
	"architectSocial/service"
	"architectSocial/templates"
	"database/sql"
	"errors"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"strconv"
)

func NewWebServer(templ *template.Template, port uint16, store sessions.Store, db *sql.DB) error {
	//http.HandleFunc("/register", createRegisterGetHandler(templ))
	http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
		var handler func(writer http.ResponseWriter, request *http.Request)
		if request.Method == http.MethodPost {
			handler = createRegisterPostHandler(service.CreateRegisterUserService(model.CreateMysqlUserModel(db)), templ)
		} else {
			handler = createRegisterGetHandler(templ)
		}

		handler(writer, request)
	})

	err := http.ListenAndServe(":"+strconv.Itoa(int(port)), nil)
	if err != nil {
		return errors.New("Failed to create a web server. Error: " + err.Error())
	}

	return nil
}

func createRegisterGetHandler(templ *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templ.ExecuteTemplate(w, "register.html", templates.RegisterData{
			PageTitle: "Register New User",
			Errors:    []string{},
		})
		if err != nil {
			showErrorPage(w, templ)
		}
	}
}

func createRegisterPostHandler(registerUser service.RegisterUserService, templ *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			showErrorPage(w, templ)
		}
		dto := service.RegisterUserDto{
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
				showErrorPage(w, templ)
			}
			return
		}
	}
}

func showErrorPage(w http.ResponseWriter, templ *template.Template) {
	_ = templ.ExecuteTemplate(w, "error.html", struct{ PageTitle string }{
		PageTitle: "Register New User",
	})
}
