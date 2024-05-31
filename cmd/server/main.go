package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/aprole/books-api/pkg/book"

	"github.com/gorilla/mux"
)

var books sync.Map

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bookSlice []book.Book
	books.Range(func(key, value interface{}) bool {
		bookSlice = append(bookSlice, value.(book.Book))
		return true
	})

	json.NewEncoder(w).Encode(bookSlice)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	b, ok := books.Load(id)
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(b)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book book.Book
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if _, ok := books.Load(book.ID); ok {
		http.Error(w, fmt.Sprintf("Book ID %q already exists", book.ID), http.StatusConflict)
		return
	}
	books.Store(book.ID, book)
	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	defer r.Body.Close()
	var updatedBook book.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if _, ok := books.Load(id); !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	books.Store(id, updatedBook)
	json.NewEncoder(w).Encode(updatedBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	book, ok := books.Load(id)
	if !ok {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	books.Delete(id)
	json.NewEncoder(w).Encode(book)
}

func main() {
	r := mux.NewRouter()

	books.Store("1", book.Book{ID: "1", Title: "1984", Author: "George Orwell"})
	books.Store("2", book.Book{ID: "2", Title: "Brave New World", Author: "Aldous Huxley"})

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	port := 8000

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Server started and listening on port %d\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %d: %v\n", port, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("Server gracefully Shut down")
}
