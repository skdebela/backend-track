package services

import (
	"errors"
	"sync"

	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	AddMember(member models.Member)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ReserveBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) ([]models.Book, error)
}

type ReservationRequest struct {
	BookID   int
	MemberID int
	Response chan<- error
}

type Library struct {
	Books			map[int]models.Book
	Members			map[int]models.Member
	Mu				sync.Mutex
	ReservationChan chan ReservationRequest
}

func NewLibrary() *Library {
	return &Library{
		Books:				make(map[int]models.Book),
		Members:			make(map[int]models.Member),
		ReservationChan:	make(chan ReservationRequest, 100),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	l.Books[book.ID] = book
}

func (l *Library) AddMember(member models.Member) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	l.Members[member.ID] = member
}

func (l *Library) RemoveBook(bookID int) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	book, exists := l.Books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	if book.Status == models.Borrowed {
		return errors.New("Cannot remove borrowed book")
	}
	delete(l.Books, bookID)
	return nil
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	
	book, bookExists := l.Books[bookID]
	if !bookExists {
        return errors.New("book not found")
    }
    if book.Status == models.Borrowed {
        return errors.New("book is already borrowed")
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
        return errors.New("member not found")
    }

	book.Status = models.Borrowed
	l.Books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	book, bookExists := l.Books[bookID]
	if !bookExists {
        return errors.New("book not found")
    }
    if book.Status == models.Available {
        return errors.New("book is not borrowed")
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
        return errors.New("member not found")
	}
	for i, b := range member.BorrowedBooks {
		if b.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			break
		}
	}
	l.Members[memberID] = member

	book.Status = models.Available 
	l.Books[bookID] = book

	return nil
}

func (l *Library) ReserveBook(bookID int, memberID int) error {
	resp := make(chan error, 1)
	req := ReservationRequest{
		BookID:   bookID,
		MemberID: memberID,
		Response: resp,
	}
	l.ReservationChan <- req
	return <-resp
}

func (l *Library) ListAvailableBooks() []models.Book {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	var available []models.Book
	for _, book := range l.Books {
		if book.Status == models.Available {
			available = append(available, book)
		}
	}

	return available
}

func (l *Library) ListBorrowedBooks(memberID int) ([]models.Book, error) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	member, exists := l.Members[memberID]
	if !exists {
		return nil, errors.New("member not found")
	}
	return member.BorrowedBooks, nil
}