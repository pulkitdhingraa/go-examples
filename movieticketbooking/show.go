package main

import (
	"time"
)

type Show struct {
	ID        string
	Movie     *Movie
	Theater   *Theater
	StartTime time.Time
	EndTime   time.Time
	Seats     map[string]*Seat
}

func NewShow(id string, movie *Movie, theater *Theater, st, et time.Time, seats map[string]*Seat) *Show {
	return &Show{
		ID:        id,
		Movie:     movie,
		Theater:   theater,
		StartTime: st,
		EndTime:   et,
		Seats:     seats,
	}
}
