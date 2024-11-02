package main

import (
	"database/sql"
	"flag"
	"fmt"
	"garlip/internal/handler"
	"garlip/internal/queries"
	"garlip/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "specify a port number on which to start the application")
	flag.Parse()
}

func main() {
	log.Println("Connecting to", os.Getenv("DB_URL"))
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	qs := queries.New(db)
	authSvc := &service.AuthService{Queries: qs}

	mux := chi.NewMux()
	apiV1 := chi.NewRouter()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/helloworld", handleHelloWorld)
	mux.Mount("/api/v1", apiV1)

	apiV1.Route("/auth", func(r chi.Router) {
		authAPI := &handler.AuthHandler{AuthSVC: authSvc}

		r.Post("/register", authAPI.Register)
		r.Post("/login", authAPI.Login)
	})

	log.Printf("Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", r.Header.Get("content-type"))
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "Hello world!")
}
