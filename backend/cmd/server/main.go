package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
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
