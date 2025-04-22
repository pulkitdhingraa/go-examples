package main

func main() {
	library := GetLibrary()

	library.AddBook("Da Vinci Code", "Dan Brown")
	library.AddBook("Inferno", "Dan Brown")
	library.AddBook("Digital Fortress", "Dan Brown")
	library.AddBook("Angels and Demons", "Dan Brown")
	library.AddBook("Origin", "Dan Brown")

	library.AddUser("PD")

	library.BorrowBook(1, 3)
	library.CheckOverdue(1)
	library.ReturnBook(1,2)
	library.ReturnBook(1,3)
}