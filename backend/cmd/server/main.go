package main

import (
	"database/sql"
	"flag"
	"garlip/internal/postgres"
	"garlip/internal/service"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "specify a port number on which to start the application")
	flag.Parse()
}

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

	log.Printf("Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleEcho(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", r.Header.Get("content-type"))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, r.Body)
}
