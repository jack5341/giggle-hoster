package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jack5341/giggle-hoster/api/handler"
	"github.com/jack5341/giggle-hoster/api/middleware"
	"github.com/jack5341/giggle-hoster/internal/database"
	"github.com/jack5341/giggle-hoster/internal/redis"
	"github.com/jack5341/giggle-hoster/internal/types"
)

func main() {
	cache, err := redis.EstablishRedisConnection()
	if err != nil {
		log.Fatal(err)
	}

	redis.CacheConn = cache

	db, err := database.EstablishDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(types.Node{}); err != nil {
		log.Fatal(err)
	}

	database.Db = db

	r := mux.NewRouter()

	authz := r.PathPrefix("/auth").Subrouter()
	authz.HandleFunc("/sigin", handler.SignIn).Methods(http.MethodPost)
	authz.HandleFunc("/sigup", handler.SignUp).Methods(http.MethodPost)
	authz.HandleFunc("/password-forgot", nil).Methods(http.MethodPost)
	authz.HandleFunc("/password-reset", nil).Methods(http.MethodPost)
	authz.HandleFunc("/email-verification", nil).Methods(http.MethodPost)

	podz := r.PathPrefix("/node").Subrouter()
	podz.Use(middleware.AuthMiddleware)
	podz.HandleFunc("/", nil).Methods(http.MethodPost)
	podz.HandleFunc("/", nil).Methods(http.MethodGet)
	podz.HandleFunc("/", nil).Methods(http.MethodDelete)

	http.Handle("/", r)

	err = http.ListenAndServe(":3000", r)
	log.Fatal(err)
}
