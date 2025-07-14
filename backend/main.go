package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	InitDB()

	r := mux.NewRouter()
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	r.HandleFunc("/notes", AuthMiddleware(CreateNoteHandler)).Methods("POST")
	r.HandleFunc("/notes", AuthMiddleware(GetNotesHandler)).Methods("GET")
	r.HandleFunc("/notes/{id:[0-9]+}", AuthMiddleware(DeleteNoteHandler)).Methods("DELETE")
	r.HandleFunc("/notes/{id:[0-9]+}", AuthMiddleware(UpdateNoteHandler)).Methods("PUT")

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", enableCors(r)))

}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
