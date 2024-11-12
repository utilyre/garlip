package main

import (
	"context"
	"errors"
	"flag"
	"garlip/internal/handler"
	"garlip/internal/queries"
	"garlip/internal/service"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/utilyre/xmate/v3"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8080", "specify a port number on which to start the application")
	flag.Parse()
}

func main() {
	log.Println("Connecting to", os.Getenv("DB_URL"))
	db, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	pgQueries := queries.New(db)
	authSvc := &service.AuthService{Queries: pgQueries}

	mux := chi.NewMux()
	apiV1 := chi.NewRouter()
	xmate.SetDefault(handleError)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/helloworld", xmate.HandleFunc(handleHelloWorld))
	mux.Mount("/api/v1", apiV1)

	apiV1.Route("/auth", func(r chi.Router) {
		authHandler := &handler.AuthHandler{AuthSVC: authSvc}

		r.Post("/register", xmate.HandleFunc(authHandler.Register))
		r.Post("/login", xmate.HandleFunc(authHandler.Login))
	})

	log.Printf("Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "Hello world!")
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if handlerErr := (handler.Error{}); errors.As(err, &handlerErr) {
		_ = xmate.WriteJSON(w, handlerErr.Status, map[string]any{
			"message": handlerErr.Message,
		})
		return
	}
	if validationErr := (service.ValidationError{}); errors.As(err, &validationErr) {
		_ = xmate.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"field":   validationErr.Field,
			"message": validationErr.Message,
		})
		return
	}

	log.Printf("%s %s failed: %v\n", r.Method, r.URL.Path, err)
	_ = xmate.WriteJSON(w, http.StatusInternalServerError, map[string]any{
		"message": "Internal Server Error",
	})
}
