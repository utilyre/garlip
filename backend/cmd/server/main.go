package main

import (
	"database/sql"
	"garlip/internal/postgres"
	"garlip/internal/service"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	queries := postgres.New(db)
	authSvc := service.Auth{Queries: queries}

	_ = authSvc

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", handleEcho)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe("", mux))
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", r.Header.Get("content-type"))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, r.Body)
}
