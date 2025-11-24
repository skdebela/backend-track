package models

type BookStatus string

const (
	Available BookStatus = "Available"
	Borrowed BookStatus = "Borrowed"
	Reserved BookStatus = "Reserved"
)

type Book struct {
	ID			int
	Title		string
	Author		string
	Status		BookStatus
	ReservedBy	int  // memberId if reserved, 0 if not
}

func NewBook(id int, title string, author string) Book {
	return Book{
		ID:			id,
		Title:		title,
		Author:		author,
		Status:		Available,
		ReservedBy: 0,
	}
}