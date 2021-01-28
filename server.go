package main

import (
	"log"
	"net/http"

	"server/handlers"
	"server/store"

	"github.com/gorilla/mux"
)

// Args args used to run the server
type Args struct {
	// postgres connection string, of the form,
	// e.g "postgres://user:password@localhost:5432/database?sslmode=disable"
	conn string
	// port for the server of the form,
	// e.g ":8080"
	port string
}

// Run run the server based on given args
func Run(args Args) error {
	// router
	router := mux.NewRouter().
		PathPrefix("/api/v1/"). // add prefix for v1 api `/api/v1/`
		Subrouter()

	st := store.NewPostgresPostStore(args.conn)
	hnd := handlers.NewPostHandler(st)
	RegisterAllRoutes(router, hnd)

	// start server
	log.Println("Starting server at port: ", args.port)
	return http.ListenAndServe(args.port, router)
}

// RegisterAllRoutes registers all routes of the api
func RegisterAllRoutes(router *mux.Router, hnd handlers.IPostHandler) {

	// set content type
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	// get posts
	router.HandleFunc("/post", hnd.Get).Methods(http.MethodGet)
	// create posts
	router.HandleFunc("/post", hnd.Create).Methods(http.MethodPost)
	// delete post
	router.HandleFunc("/post", hnd.Delete).Methods(http.MethodDelete)

	// update post details
	// router.HandleFunc("/post/details", hnd.UpdateDetails).Methods(http.MethodPut)

	// list posts
	router.HandleFunc("/posts", hnd.List).Methods(http.MethodGet)

	// list Reddit top posts
	router.HandleFunc("/reddit", hnd.RedditList).Methods(http.MethodGet)
}
