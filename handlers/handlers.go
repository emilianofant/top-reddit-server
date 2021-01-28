package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/errors"
	"server/objects"
	"server/store"
	"time"
)

// IPostHandler is implement all the handlers
type IPostHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	RedditList(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	UpdateViewed(w http.ResponseWriter, r *http.Request)
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
	// limit
	limit, err := IntFromString(w, values.Get("limit"))
	if err != nil {
		return
	}
	// list Posts
	list, err := h.store.List(r.Context(), &objects.ListRequest{
		Limit: limit,
		After: after,
	})
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.PostResponseWrapper{Posts: list})
}

func (h *handler) RedditList(w http.ResponseWriter, r *http.Request) {
	var resList []*objects.Post

	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", "https://www.reddit.com/r/gaming/top.json", nil)
	req.Header.Set("User-Agent", "reddit-top/1.0")

	res, err := httpClient.Do(req)

	// check for response error
	if err != nil {
		WriteError(w, errors.ErrBadRequest)
		return
	}

	data, _ := ioutil.ReadAll(res.Body)

	// close response body
	defer res.Body.Close()

	type response struct {
		Data struct {
			Children []struct {
				Data objects.Post
			} `json:"children"`
		} `json:"data"`
	}

	var redditData response

	jsonErr := json.Unmarshal(data, &redditData)
	if jsonErr != nil {
		WriteError(w, jsonErr)
	}

	// for {key}, {value} := range {list}
	for _, PostWrapper := range redditData.Data.Children {
		post := PostWrapper.Data
		resList = append(resList, &post)
	}

	WriteResponse(w, &objects.PostResponseWrapper{Posts: resList})
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

func (h *handler) UpdateViewed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, errors.ErrUnprocessableEntity)
		return
	}
	req := &objects.UpdateRequest{}
	if Unmarshal(w, data, req) != nil {
		return
	}

	// check if Post exist
	if _, err := h.store.Get(r.Context(), &objects.GetRequest{ID: req.ID}); err != nil {
		WriteError(w, err)
		return
	}

	if err = h.store.UpdateViewed(r.Context(), req); err != nil {
		WriteError(w, err)
		return
	}
	WriteResponse(w, &objects.PostResponseWrapper{})
}

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
