package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	hash, _ := HashPassword(user.Password)
	err := DB.QueryRow(r.Context(), "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id", user.Email, hash).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	token, _ := GenerateJWT(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input User
	json.NewDecoder(r.Body).Decode(&input)

	var user User
	err := DB.QueryRow(r.Context(), "SELECT id, password FROM users WHERE email=$1", input.Email).Scan(&user.ID, &user.Password)
	if err != nil || CheckPasswordHash(input.Password, user.Password) != nil {
		http.Error(w, "pswds aint match honey", http.StatusUnauthorized)
		return
	}
	token, _ := GenerateJWT(user.ID)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := ParseJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// r.Header.Set("X-User-ID", string(rune(userId)))
		r.Header.Set("X-User-ID", strconv.Itoa(userID))
		next(w, r)
	}
}

func GetUserID(r *http.Request) int {
	uid := r.Header.Get("X-User-ID")
	id, _ := strconv.Atoi(uid)
	return id
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	var note Note
	json.NewDecoder(r.Body).Decode(&note)

	err := DB.QueryRow(r.Context(), "INSERT INTO notes (user_id, title, body) VALUES ($1, $2, $3) RETURNING id",
		userID, note.Title, note.Body).Scan(&note.ID)

	if err != nil {
		http.Error(w, "Error creating note", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(note)
}

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	rows, err := DB.Query(r.Context(), "SELECT id, title, body FROM notes WHERE user_id=$1", userID)

	if err != nil {
		http.Error(w, "Error while fetching notes", http.StatusInternalServerError)
		return
	}

	var notes []Note
	for rows.Next() {
		var note Note
		rows.Scan(&note.ID, &note.Title, &note.Body)
		notes = append(notes, note)
	}

	json.NewEncoder(w).Encode(notes)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	noteIDStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	tag, err := DB.Exec(r.Context(), "DELETE FROM notes WHERE id=$1 AND user_id=$2", noteID, userID)
	if err != nil {
		http.Error(w, "Error deleting note", http.StatusInternalServerError)
		return
	}

	if tag.RowsAffected() == 0 {
		http.Error(w, "Note not found or not yours", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := GetUserID(r)
	noteIDStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tag, err := DB.Exec(r.Context(),
		"UPDATE notes SET title=$1, body=$2 WHERE id=$3 AND user_id=$4",
		note.Title, note.Body, noteID, userID,
	)
	if err != nil {
		http.Error(w, "Error updating note", http.StatusInternalServerError)
		return
	}

	if tag.RowsAffected() == 0 {
		http.Error(w, "Note not found or not yours", http.StatusNotFound)
		return
	}

	note.ID = noteID
	json.NewEncoder(w).Encode(note)
}
