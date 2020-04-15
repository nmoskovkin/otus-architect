package helpers

import (
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
)

type SessionWrapper interface {
	SetRegistrationId(id string, r *http.Request, w http.ResponseWriter) error
	SetAuthenticated(id string, r *http.Request, w http.ResponseWriter) error
	GetRegistrationId(r *http.Request) (string, error)
	IsAuthenticated(r *http.Request) (bool, string, error)
}

type GorillaSessionWrapper struct {
	store sessions.Store
}

func NewGorillaSessionWrapper(store sessions.Store) *GorillaSessionWrapper {
	return &GorillaSessionWrapper{store: store}
}

func (wrapper *GorillaSessionWrapper) SetRegistrationId(id string, r *http.Request, w http.ResponseWriter) error {
	session, _ := wrapper.store.Get(r, "user-session")
	session.Values["registration-id"] = id
	err := session.Save(r, w)
	if err != nil {
		return errors.New("failed to save register id, error: " + err.Error())
	}

	return nil
}

func (wrapper *GorillaSessionWrapper) SetAuthenticated(id string, r *http.Request, w http.ResponseWriter) error {
	session, _ := wrapper.store.Get(r, "user-session")
	session.Values["is-authenticated"] = true
	session.Values["auth-id"] = id
	err := session.Save(r, w)
	if err != nil {
		return errors.New("failed to save \"is-authenticated\", error: " + err.Error())
	}

	return nil
}

func (wrapper *GorillaSessionWrapper) GetRegistrationId(r *http.Request) (string, error) {
	session, err := wrapper.store.Get(r, "user-session")
	if err != nil {
		return "", errors.New("failed to get registration id, error: " + err.Error())
	}
	id, ok := session.Values["registration-id"]
	if !ok {
		return "", errors.New("failed to get registration id")
	}
	v, ok := id.(string)
	if !ok {
		return "", errors.New("failed to get registration id")
	}
	return v, nil
}

func (wrapper *GorillaSessionWrapper) IsAuthenticated(r *http.Request) (bool, string, error) {
	session, err := wrapper.store.Get(r, "user-session")
	if err != nil {
		return false, "", errors.New("failed to get \"is-authenticated\", error: " + err.Error())
	}
	isAuth, ok := session.Values["is-authenticated"]
	if !ok {
		return false, "", errors.New("failed to get \"is-authenticated\"")
	}
	isAuthV, ok := isAuth.(bool)
	if !ok {
		return false, "", errors.New("failed to get \"is-authenticated\"")
	}
	id, ok := session.Values["auth-id"]
	if !ok {
		return false, "", errors.New("failed to get auth id")
	}
	idV, ok := id.(string)

	return isAuthV, idV, nil
}
