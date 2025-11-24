package concurrency

import (
	"errors"
	"fmt"
	"time"

	"library_management/models"
	"library_management/services"
)

func StartWorkers(l *services.Library, numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go worker(l)
	}
}

func worker(l *services.Library) {
	for req := range l.ReservationChan {
		processReservation(req, l)
	}
}

func processReservation(req services.ReservationRequest, l *services.Library) {
	l.Mu.Lock()
	defer l.Mu.Unlock()

	book, ok := l.Books[req.BookID]
	if !ok {
		req.Response <- errors.New("book not found")
		return
	}

	if book.Status != models.Available {
		req.Response <- errors.New("book not available")
		return
	}

	book.Status = models.Reserved
	book.ReservedBy = req.MemberID
	l.Books[req.BookID] = book

	req.Response <- nil

	// Start auto-cancellation timer
	go func(bid, mid int) {
		time.Sleep(5 * time.Second)
		l.Mu.Lock()
		defer l.Mu.Unlock()
		if b, ok := l.Books[bid]; ok && b.Status == models.Reserved && b.ReservedBy == mid {
			b.Status = models.Available
			b.ReservedBy = 0
			l.Books[bid] = b
			fmt.Printf("Auto-unreserved book ID %d after timeout\n", bid)
		}
	}(req.BookID, req.MemberID)
}