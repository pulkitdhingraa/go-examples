package main

type Movie struct {
	ID              string
	Title           string
	DurationMinutes int
}

func NewMovie(id, title string, duration int) *Movie {
	return &Movie{
		ID: id,
		Title: title,
		DurationMinutes: duration,
	}
}