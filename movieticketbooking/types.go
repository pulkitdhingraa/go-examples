package main

type SeatType int
type SeatStatus int
type BookingStatus int

const (
	Gold SeatType = iota
	Platinum
)

const (
	Available SeatStatus = iota
	Booked
)

const (
	Pending BookingStatus = iota
	Confirmed
	Cancelled
)
