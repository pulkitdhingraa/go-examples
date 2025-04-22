package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Library struct {
	Books map[int]*Book
	Users map[int]*User
	mu    sync.RWMutex
}

var (
	once sync.Once
	instance *Library
)

func GetLibrary() *Library{
	once.Do(func() {
		instance = &Library{
			Books: make(map[int]*Book),
			Users: make(map[int]*User),
		}
	})
	return instance
}

func (l *Library) AddBook(title, author string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	bookID := len(l.Books)+1
	book := Book{ID: bookID, Title: title, Author: author}
	l.Books[bookID] = &book
	fmt.Printf("Added Book: %s by %s\n", title, author)
}

func (l *Library) AddUser(name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	userID := len(l.Users)+1
	user := User{ID: userID, Name: name, BorrowedBooks: make(map[int]time.Time)}
	l.Users[userID] = &user
	fmt.Printf("Added user: %s\n", name)
}

func (l *Library) BorrowBook(userID, bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Find the user with userID
	user := l.getUserByID(userID)
	// Find the book with bookID
	book := l.getBookByID(bookID)

	// if user or book doesn't exist
	if user == nil || book == nil {
		fmt.Printf("User or Book not found\n")
		return
	}

	// if book is already borrowed
	if book.IsBorrowed {
		fmt.Printf("Book already borrowed\n")
		return
	}

	book.IsBorrowed = true
	user.BorrowedBooks[bookID] = time.Now()
	fmt.Printf("%s borrowed %s\n", user.Name, book.Title)
}

func (l *Library) ReturnBook(userID, bookID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	user := l.getUserByID(userID)
	book := l.getBookByID(bookID)

	if user == nil || book == nil {
		fmt.Printf("User or Book not found\n")
		return
	}

	if _, ok := user.BorrowedBooks[bookID]; !ok {
		fmt.Printf("Book %s wasn't borrowed by the user %s\n", book.Title, user.Name)
		return
	}

	delete(user.BorrowedBooks, bookID)
	book.IsBorrowed = false
	fmt.Printf("%s returned %s\n", user.Name, book.Title)
}

func (l *Library) CheckOverdue(userID int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	user := l.getUserByID(userID)
	if user == nil {
		fmt.Println("User not found in the system")
		return
	}

	for bookID, borrowedTime := range user.BorrowedBooks {
		days := int(time.Since(borrowedTime).Hours()/24)
		if days > 7 {
			fmt.Printf("Book ID %d is overdue by %d days\n", bookID, days-7)
		} else {
			fmt.Printf("No overdue for the user %s\n", user.Name)
		}
	}
}

func (l *Library) getBookByID(bookID int) *Book{
	return l.Books[bookID]
}

func (l *Library) getUserByID(userID int) *User {
	return l.Users[userID]
}

func (l *Library) SearchBooks(keyword string) []*Book {
	keyword = strings.ToLower(keyword)
	match := make([]*Book, 0)

	for _, book := range l.Books {
		if strings.Contains(strings.ToLower(book.Title), keyword) || strings.Contains(strings.ToLower(book.Author), keyword) {
			match = append(match, book)
		}
	}
	return match
}