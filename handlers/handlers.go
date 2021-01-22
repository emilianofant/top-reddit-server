package handlers

import (
	"net/http"
)

// IPostHandler is implement all the handlers
type IPostHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateDetails(w http.ResponseWriter, r *http.Request)
	Cancel(w http.ResponseWriter, r *http.Request)
	// Reschedule(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
}

// NewPostHandler return current IPostHandler implementation
func NewPostHandler() IPostHandler {
	return &handler{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *handler) Cancel(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
