package controller

import (
	"html/template"
	"net/http"
)

type ErrorReturningHandlerFunc func(writer http.ResponseWriter, request *http.Request) error
type ErrorHandlerFunc func(writer http.ResponseWriter, request *http.Request, clientError ClientError)

type Handler struct {
	rootHandler  ErrorReturningHandlerFunc
	errorHandler ErrorHandlerFunc
}

func NewHandler(rootHandler ErrorReturningHandlerFunc, templ *template.Template) *Handler {
	return &Handler{
		rootHandler:  rootHandler,
		errorHandler: CreateErrorHandler(templ),
	}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := handler.rootHandler(writer, request)
	if err == nil {
		return
	}

	clientError, ok := err.(ClientError)
	if !ok {
		writer.WriteHeader(500)
		return
	}
	handler.errorHandler(writer, request, clientError)
}
