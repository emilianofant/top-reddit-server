package objects

import (
	"encoding/json"
	"net/http"
)

// MaxListLimit maximum listting
const MaxListLimit = 200

// GetRequest for retrieving single Post
type GetRequest struct {
	ID string `json:"id"`
}

// ListRequest for retrieving list of Posts
type ListRequest struct {
	Limit int    `json:"limit"`
	After string `json:"after"`
}

// CreateRequest for creating a new Fav Post
type CreateRequest struct {
	Post *Post `json:"post"`
}

// DeleteRequest to delete a Post
type DeleteRequest struct {
	ID string `json:"id"`
}

// PostResponseWrapper response of any Post request
type PostResponseWrapper struct {
	Post  *Post   `json:"post,omitempty"`
	Posts []*Post `json:"posts,omitempty"`
	Code  int     `json:"-"`
}

// JSON convert PostResponseWrapper in json
func (e *PostResponseWrapper) JSON() []byte {
	if e == nil {
		return []byte("{}")
	}
	res, _ := json.Marshal(e)
	return res
}

// StatusCode return status code
func (e *PostResponseWrapper) StatusCode() int {
	if e == nil || e.Code == 0 {
		return http.StatusOK
	}
	return e.Code
}
