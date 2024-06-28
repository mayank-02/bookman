package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/mayank-02/bookman/internal/db"
	"github.com/mayank-02/bookman/internal/models"

	"github.com/gorilla/mux"
)

const (
	APIVersion      = "v1"
	BooksPath       = "/api/" + APIVersion + "/books"
	CollectionsPath = "/api/" + APIVersion + "/collections"
)

func RegisterHandlers(r *mux.Router, db *db.DB) {
	r.HandleFunc(BooksPath, getBooks(db)).Methods("GET")
	r.HandleFunc(BooksPath, createBook(db)).Methods("POST")
	r.HandleFunc(BooksPath+"/{id}", getBook(db)).Methods("GET")
	r.HandleFunc(BooksPath+"/{id}", updateBook(db)).Methods("PUT")
	r.HandleFunc(BooksPath+"/{id}", deleteBook(db)).Methods("DELETE")
	r.HandleFunc(CollectionsPath, getCollections(db)).Methods("GET")
	r.HandleFunc(CollectionsPath, createCollection(db)).Methods("POST")
	r.HandleFunc(CollectionsPath+"/{id}", getCollection(db)).Methods("GET")
	r.HandleFunc(CollectionsPath+"/{id}", updateCollection(db)).Methods("PUT")
	r.HandleFunc(CollectionsPath+"/{id}", deleteCollection(db)).Methods("DELETE")
	r.HandleFunc(CollectionsPath+"/{id}/books/{bookId}", addBookToCollection(db)).Methods("POST")
	r.HandleFunc(CollectionsPath+"/{id}/books/{bookId}", removeBookFromCollection(db)).Methods("DELETE")
}

func getBooks(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		author := r.URL.Query().Get("author")
		genre := r.URL.Query().Get("genre")
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		books, err := db.GetBooks(author, genre, from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(books)
	}
}

func getBook(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		book, err := db.GetBook(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

func createBook(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validation checks
		if book.Title == "" || book.Author == "" || book.PublishedDate == "" {
			http.Error(w, "Title, author, and published date are required", http.StatusBadRequest)
			return
		}

		if result, err := time.Parse("2006-01-02", book.PublishedDate); err != nil || result.IsZero() {
			http.Error(w, "Invalid published date format, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		id, err := db.CreateBook(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		book.ID = id
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func updateBook(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		var book models.Book
		json.NewDecoder(r.Body).Decode(&book)
		book.ID = id

		// Validation checks
		if book.Title == "" || book.Author == "" || book.PublishedDate == "" {
			http.Error(w, "Title, author, and published date are required", http.StatusBadRequest)
			return
		}

		if result, err := time.Parse("2006-01-02", book.PublishedDate); err != nil || result.IsZero() {
			http.Error(w, "Invalid published date format, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		err := db.UpdateBook(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

func deleteBook(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		err := db.DeleteBook(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func getCollections(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collections, err := db.GetCollections()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(collections)
	}
}

func getCollection(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		collection, err := db.GetCollection(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(collection)
	}
}

func createCollection(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var collection models.Collection
		json.NewDecoder(r.Body).Decode(&collection)

		// Validation checks
		if collection.Name == "" {
			http.Error(w, "Collection name is required", http.StatusBadRequest)
			return
		}

		id, err := db.CreateCollection(collection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		collection.ID = id
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(collection)
	}
}

func updateCollection(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		var collection models.Collection
		json.NewDecoder(r.Body).Decode(&collection)
		collection.ID = id

		// Validation checks
		if collection.Name == "" {
			http.Error(w, "Collection name is required", http.StatusBadRequest)
			return
		}

		err := db.UpdateCollection(collection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(collection)
	}
}

func deleteCollection(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(mux.Vars(r)["id"])
		err := db.DeleteCollection(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
func addBookToCollection(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collectionID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "invalid collection ID", http.StatusBadRequest)
			return
		}
		bookID, err := strconv.Atoi(vars["bookId"])
		if err != nil {
			http.Error(w, "invalid book ID", http.StatusBadRequest)
			return
		}

		// Check if the book is already in the collection
		inCollection, err := db.IsBookInCollection(collectionID, bookID)
		if err != nil {
			http.Error(w, "error checking book in collection: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if inCollection {
			http.Error(w, "book is already in the collection", http.StatusConflict)
			return
		}

		// Add the book to the collection
		err = db.AddBookToCollection(collectionID, bookID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func removeBookFromCollection(db *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		collectionID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "invalid collection ID", http.StatusBadRequest)
			return
		}
		bookID, err := strconv.Atoi(vars["bookId"])
		if err != nil {
			http.Error(w, "invalid book ID", http.StatusBadRequest)
			return
		}

		// Check if the book is in the collection
		inCollection, err := db.IsBookInCollection(collectionID, bookID)
		if err != nil {
			http.Error(w, "error checking book in collection: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !inCollection {
			http.Error(w, "book is not in the collection", http.StatusNotFound)
			return
		}

		// Remove the book from the collection
		err = db.RemoveBookFromCollection(collectionID, bookID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
