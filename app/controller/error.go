package controller

import (
	"architectSocial/app/templates"
	"html/template"
	"net/http"
)

type ClientError interface {
	Error() string
}

type HTTPError struct {
	Cause  error
	Detail string
	Status int
}

func (e *HTTPError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}
	return e.Detail + " : " + e.Cause.Error()
}

func NewHTTPError(err error, status int, detail string) error {
	return &HTTPError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func CreateErrorHandlerFunc(templ *template.Template) ErrorHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, clientError ClientError) {
		_ = templ.ExecuteTemplate(writer, "error.html", templates.ErrorData{
			PageTitle: "Error",
			Error:     clientError.Error(),
		})
	}
}
