package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "The Sorcerer's Stone", Author: "J.K. Rowling", Quantity: 1},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 2},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 2},
	{ID: "4", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
}

func getBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	// c.Param gets path parameter
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.JSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "All books with this title are currently checked out"})
		return
	}

	book.Quantity -= 1
	c.JSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.JSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/books/checkout", checkoutBook)
	router.PATCH("/books/return", returnBook)
	router.Run("localhost:8080")
}
