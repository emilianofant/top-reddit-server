package handlers

import (
	"io/ioutil"
	"net/http"
	"server/errors"
	"server/objects"
	"server/store"
)

// IPostHandler is implement all the handlers
type IPostHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	// UpdateDetails(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	store store.IPostStore
}

func NewPostHandler(store store.IPostStore) IPostHandler {
	return &handler{store: store}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, errors.ErrValidPostIDIsRequired)
		return
	}
	post, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.PostResponseWrapper{Post: post})
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	// after
	after := values.Get("after")
	// name
	name := values.Get("name")
	// limit
	limit, err := IntFromString(w, values.Get("limit"))
	if err != nil {
		return
	}
	// list Posts
	list, err := h.store.List(r.Context(), &objects.ListRequest{
		Limit: limit,
		After: after,
		Name:  name,
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.PostResponseWrapper{Posts: list})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	post := &objects.Post{}
	if Unmarshal(w, data, post) != nil {
		return
	}
	if err = h.store.Create(r.Context(), &objects.CreateRequest{Post: post}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.PostResponseWrapper{Post: post})
}

// func (h *handler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
// 	data, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		WriteError(w, errors.ErrUnprocessableEntity)
// 		return
// 	}
// 	req := &objects.UpdateDetailsRequest{}
// 	if Unmarshal(w, data, req) != nil {
// 		return
// 	}

// 	// check if Post exist
// 	if _, err := h.store.Get(r.Context(), &objects.GetRequest{Id: req.Id}); err != nil {
// 		WriteError(w, err)
// 		return
// 	}

// 	if err = h.store.UpdateDetails(r.Context(), req); err != nil {
// 		WriteError(w, err)
// 		return
// 	}
// 	WriteResponse(w, &objects.PostResponseWrapper{})
// }

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteError(w, errors.ErrValidPostIDIsRequired)
		return
	}

	// check if Post exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}

	if err := h.store.Delete(r.Context(), &objects.DeleteRequest{ID: id}); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.PostResponseWrapper{})
}
