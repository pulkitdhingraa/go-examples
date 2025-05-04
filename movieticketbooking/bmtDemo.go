package main

import (
	"fmt"
	"time"
)

func main() {
	bmtSystem := NewBMTSystem()

	movie1 := NewMovie("M1", "Movie 1", 120)
	movie2 := NewMovie("M2", "Movie 2", 180)
	bmtSystem.AddMovie(movie1)
	bmtSystem.AddMovie(movie2)

	theater1 := NewTheater("T1", "Theater 1", "Hyd")
	theater2 := NewTheater("T2", "Theater 2", "Del")
	bmtSystem.AddTheater(theater1)
	bmtSystem.AddTheater(theater2)

	show1 := NewShow("S1", movie1, theater1, time.Now(), time.Now().Add(time.Duration(movie1.DurationMinutes)*time.Minute), CreateSeatsUtility(8,8))

	show2 := NewShow("S2", movie2, theater2, time.Now(), time.Now().Add(time.Duration(movie2.DurationMinutes)*time.Minute), CreateSeatsUtility(10,10))

	bmtSystem.AddShow(show1)
	bmtSystem.AddShow(show2)

	user := NewUser("U1", "Joe", "joe@xyz.com")
	selectedSeats := []*Seat{
		show1.Seats["2-5"],
		show1.Seats["1-6"],
	}

	booking, err := bmtSystem.BookTickets(user, show1, selectedSeats)
	if err != nil {
		fmt.Printf("booking failed: %v\n", err)
		return
	}

	fmt.Printf("Booking Successful. Booking ID: %s\n", booking.ID)

	if err := bmtSystem.ConfirmBooking(booking.ID); err != nil {
		fmt.Printf("Failed to confirm booking: %v\n", err)
	}

	fmt.Println("Booking confirmed")

	fmt.Printf("Total Price of booking: %.2f", booking.TotalPrice)
}