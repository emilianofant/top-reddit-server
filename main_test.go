package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"server/handlers"
	"server/objects"
	"server/store"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	router    *mux.Router
	flushAll  func(t *testing.T)
	createOne func(t *testing.T, name string) *objects.Post
	getOne    func(t *testing.T, id string, wantErr bool) *objects.Post
)

func TestMain(t *testing.M) {
	log.Println("Registering")

	conn := "postgres://user:password@localhost:5432/db?sslmode=disable"
	if c := os.Getenv("DB_CONN"); c != "" {
		conn = c
	}

	router = mux.NewRouter().PathPrefix("/api/v1/").Subrouter()
	st := store.NewPostgresPostStore(conn)
	hnd := handlers.NewPostHandler(st)
	RegisterAllRoutes(router, hnd)

	flushAll = func(t *testing.T) {
		db, err := gorm.Open(postgres.Open(conn), nil)
		if err != nil {
			t.Fatal(err)
		}
		db.Delete(&objects.Post{}, "1=1")
	}

	createOne = func(t *testing.T, title string) *objects.Post {
		post := &objects.Post{
			Title:     title,
			Author:    "TheoneAndonly",
			EntryDate: time.Now().UTC().String(),
		}
		err := st.Create(context.TODO(), &objects.CreateRequest{Post: post})
		if err != nil {
			t.Fatal(err)
		}
		return post
	}
	getOne = func(t *testing.T, id string, wantErr bool) *objects.Post {
		post, err := st.Get(context.TODO(), &objects.GetRequest{ID: id})
		if err != nil && wantErr {
			t.Fatal(err)
		}
		return post
	}

	log.Println("Starting")
	os.Exit(t.Run())
}

func Do(req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestListEndpoint(t *testing.T) {
	flushAll(t)
	tests := []struct {
		name    string
		code    int
		setup   func(t *testing.T) *http.Request
		listLen int
	}{
		{
			name: "Zero",
			setup: func(t *testing.T) *http.Request {
				flushAll(t)
				req, err := http.NewRequest(http.MethodGet, "/api/v1/posts", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 0,
		},
		{
			name: "All",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "One")
				_ = createOne(t, "Two")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/posts", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "Limited",
			setup: func(t *testing.T) *http.Request {
				_ = createOne(t, "Three")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/posts?limit=2", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 2,
		},
		{
			name: "After",
			setup: func(t *testing.T) *http.Request {
				post := createOne(t, "Four")
				_ = createOne(t, "Five")
				req, err := http.NewRequest(http.MethodGet, "/api/v1/posts?after="+post.ID, nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 1,
		},
		{
			name: "Name",
			setup: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/api/v1/posts?name=e", nil)
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			code:    http.StatusOK,
			listLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Do(tt.setup(t))
			got := &objects.PostResponseWrapper{}
			assert.Equal(t, tt.code, w.Code)
			assert.Nil(t, json.Unmarshal(w.Body.Bytes(), got))
			assert.Equal(t, len(got.Posts), tt.listLen)
		})
	}
}
