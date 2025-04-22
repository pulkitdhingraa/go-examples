package main

import "time"

type User struct {
	ID            int
	Name          string
	BorrowedBooks map[int]time.Time
}