package main

import (
	"database/sql"
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
	_ "github.com/lib/pq"
	"github.com/utilyre/xmate/v2"
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
	eh := xmate.ErrorHandler(handleError)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Get("/helloworld", eh.HandleFunc(handleHelloWorld))
	mux.Mount("/api/v1", apiV1)

	apiV1.Route("/auth", func(r chi.Router) {
		authAPI := &handler.AuthHandler{AuthSVC: authSvc}

		r.Post("/register", eh.HandleFunc(authAPI.Register))
		r.Post("/login", eh.HandleFunc(authAPI.Login))
	})

	log.Printf("Listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleHelloWorld(w http.ResponseWriter, r *http.Request) error {
	return xmate.WriteText(w, http.StatusOK, "Hello world!")
}

func handleError(w http.ResponseWriter, r *http.Request) {
	err := r.Context().Value(xmate.KeyError).(error)

	if httpErr := (xmate.HTTPError{}); errors.As(err, &httpErr) {
		_ = xmate.WriteJSON(w, httpErr.Code, map[string]any{
			"message": httpErr.Message,
		})
		return
	}
	if validationErr := (service.ValidationError{}); errors.As(err, &validationErr) {
		_ = xmate.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"field":   validationErr.Param,
			"message": validationErr.Msg,
		})
		return
	}

	log.Printf("%s %s failed: %v\n", r.Method, r.URL.Path, err)
	_ = xmate.WriteJSON(w, http.StatusInternalServerError, map[string]any{
		"message": "Internal Server Error",
	})
}
