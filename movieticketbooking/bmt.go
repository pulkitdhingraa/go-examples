package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type BMTSystem struct {
	movies       []*Movie
	theaters     []*Theater
	shows        map[string]*Show
	bookings     map[string]*Booking
	mu           sync.RWMutex
	bookingCount int64
}

var (
	instance *BMTSystem
	once     sync.Once
)

func NewBMTSystem() *BMTSystem{
	once.Do(func() {
		instance = &BMTSystem{
			movies: make([]*Movie, 0),
			theaters: make([]*Theater, 0),
			shows: make(map[string]*Show),
			bookings: make(map[string]*Booking),
		}
	})

	return instance
}

func (bs *BMTSystem) AddMovie(movie *Movie){
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.movies = append(bs.movies, movie)
}

func (bs *BMTSystem) AddTheater(theater *Theater){
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.theaters = append(bs.theaters, theater)
}

func (bs *BMTSystem) AddShow(show *Show){
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.shows[show.ID] = show
}

func (bs *BMTSystem) GetShow(showID string) *Show{
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return bs.shows[showID]
}

func (bs *BMTSystem) BookTickets(user *User, show *Show, selectedSeats []*Seat) (*Booking, error){
	bs.mu.Lock()
	defer bs.mu.Unlock()

	for _, seat := range selectedSeats {
		showSeat, exists := show.Seats[seat.ID]
		if !exists || showSeat.GetStatus() != Available {
			return nil, fmt.Errorf("seat %s is not available", seat.ID)
		}
	}
	
	var totalPrice float64

	for _, seat := range selectedSeats{
		show.Seats[seat.ID].SetStatus(Booked)
		totalPrice += seat.GetPrice()
	}
	
	bookingID := bs.generateNewBookingID()
	booking := NewBooking(bookingID, user, show, selectedSeats, totalPrice, Pending)
	bs.bookings[bookingID] = booking

	return booking, nil 
}

func (bs *BMTSystem) ConfirmBooking(bookingID string) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	booking, exists := bs.bookings[bookingID] 
	if !exists {
		return fmt.Errorf("booking not found")
	}
	if booking.GetBookingStatus() != Pending {
		return fmt.Errorf("booking status not pending")
	}

	booking.SetBookingStatus(Confirmed)
	return nil
}

func (bs *BMTSystem) CancelBooking(bookingID string) error{
	bs.mu.Lock()
	defer bs.mu.Unlock()

	booking, exists := bs.bookings[bookingID]
	if !exists {
		return fmt.Errorf("booking not found")
	}
	if booking.GetBookingStatus() == Cancelled {
		return fmt.Errorf("Booking already cancelled")
	}
	
	booking.SetBookingStatus(Cancelled)

	for _, seat := range booking.Seats {
		booking.Show.Seats[seat.ID].SetStatus(Available)
	}

	return nil
}

func (bs *BMTSystem) generateNewBookingID() string{
	count := atomic.AddInt64(&bs.bookingCount, 1)
	return fmt.Sprintf("BMT%s%06d", time.Now().Format("20250405180530"), count)
}

func CreateSeatsUtility(rows, cols int) map[string]*Seat {
	seats := make(map[string]*Seat)
	for row:=1;row<rows;row++{
		for col:=1;col<cols;col++{
			seatid := fmt.Sprintf("%d-%d",row,col)
			seatType := Gold
			price := 300.0

			if row < 2{
				seatType = Platinum
				price = 350.0
			}

			seats[seatid] = NewSeat(seatid, row, col, seatType, price, Available)
		}
	}
	return seats
}