package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"library_management/models"
	"library_management/services"
)

type CLIController struct {
	lib services.LibraryManager
}

func NewCLIController(lib services.LibraryManager) *CLIController {
	return &CLIController{
		lib: lib,
	}
}

func (c *CLIController) Start() {
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n===== LIBRARY MANAGEMENT =====")
		fmt.Println("1. Add Book")
		fmt.Println("2. Add Member")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Reserve Book")
		fmt.Println("8. Simulate Concurretn Reservations")
		fmt.Println("9. Exit")
		fmt.Print("Select an option: ")

		reader.Scan()
		choice := reader.Text()

		switch choice {
		case "1":
			c.handleAddBook(reader)
		case "2":
			c.handleAddMember(reader)
		case "3":
			c.handleBorrowBook(reader)
		case "4":
			c.handleReturnBook(reader)
		case "5":
			c.handleListAvailable()
		case "6":
			c.handleListBorrowed(reader)
		case "7":
			c.handleReserveBook(reader)
		case "8":
			c.handleSimulateConcurrentReservations(reader)
		case "9":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Try again.")
		}
	}
}

func (c *CLIController) handleAddBook(reader *bufio.Scanner) {
	fmt.Print("Enter Book ID: ")
	id := readInt(reader)

	fmt.Print("Enter Title: ")
	reader.Scan()
	title := reader.Text()

	fmt.Print("Enter Author: ")
	reader.Scan()
	author := reader.Text()

	book := models.NewBook(id, title, author)
	c.lib.AddBook(book)
	fmt.Println("Book added successfully.")
}

func (c *CLIController) handleAddMember(reader *bufio.Scanner) {
	lib := c.lib.(*services.Library) 

	fmt.Print("Enter Member ID: ")
	id := readInt(reader)

	fmt.Print("Enter Name: ")
	reader.Scan()
	name := reader.Text()

	member := models.NewMember(id, name)
	lib.AddMember(member)

	fmt.Println("Member added successfully.")
}

func (c *CLIController) handleBorrowBook(reader *bufio.Scanner) {
	fmt.Print("Enter Book ID: ")
	bookID := readInt(reader)

	fmt.Print("Enter Member ID: ")
	memberID := readInt(reader)

	err := c.lib.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Book borrowed successfully.")
}

func (c *CLIController) handleReturnBook(reader *bufio.Scanner) {
	fmt.Print("Enter Book ID: ")
	bookID := readInt(reader)

	fmt.Print("Enter Member ID: ")
	memberID := readInt(reader)

	err := c.lib.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Book returned successfully.")
}

func (c *CLIController) handleListAvailable() {
	books := c.lib.ListAvailableBooks()
	fmt.Println("\nAvailable Books:")
	for _, b := range books {
		fmt.Printf("ID: %d | %s by %s\n", b.ID, b.Title, b.Author)
	}
}

func (c *CLIController) handleListBorrowed(reader *bufio.Scanner) {
	fmt.Print("Enter Member ID: ")
	memberID := readInt(reader)

	books, err := c.lib.ListBorrowedBooks(memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\nBorrowed Books (Member %d):\n", memberID)
	for _, b := range books {
		fmt.Printf("ID: %d | %s by %s\n", b.ID, b.Title, b.Author)
	}
}


func (c *CLIController) handleReserveBook(reader *bufio.Scanner) {
	fmt.Print("Enter Book ID: ")
	bookID := readInt(reader)

	fmt.Print("Enter Member ID: ")
	memberID := readInt(reader)

	err := c.lib.ReserveBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Book reserved successfully, You have 5 seconds to borrow it.")
}

func (c *CLIController) handleSimulateConcurrentReservations(reader *bufio.Scanner) {
	fmt.Print("Enter book ID for simulation: ")
	bookID := readInt(reader)

	fmt.Println("Simulating 5 concurrent reservation attempts by members 1-5")
	for i := 1; i <= 5; i++ {
		memberID := i
		go func(m int) {
			err := c.lib.ReserveBook(bookID, m)
			if err != nil {
				fmt.Printf("Member %d: %v\n", m, err)
			} else {
				fmt.Printf("Member %d reserved book %d successfully\n", m, bookID)
			}
		}(memberID)
	}
	fmt.Println("Launched concurrent simulations. Check console for async outputs.")
}


func readInt(reader *bufio.Scanner) int {
	for {
		reader.Scan()
		input := strings.TrimSpace(reader.Text())
		val, err := strconv.Atoi(input)
		if err != nil {
			fmt.Print("Invalid number. Try again: ")
			continue
		}
		return val
	}
}
