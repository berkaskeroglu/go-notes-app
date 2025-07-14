package main

type User struct {
	ID       int
	Email    string
	Password string
}

type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
