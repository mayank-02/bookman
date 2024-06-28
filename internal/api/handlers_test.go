package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mayank-02/bookman/internal/db"
	"github.com/mayank-02/bookman/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter(db *db.DB) *mux.Router {
	r := mux.NewRouter()
	RegisterHandlers(r, db)
	return r
}

func setupTestDB(t *testing.T) *db.DB {
	t.Helper()
	db, err := db.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		published_date TEXT NOT NULL,
		edition TEXT,
		description TEXT,
		genre TEXT,
		created_at TEXT DEFAULT (datetime('now')),
		updated_at TEXT DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at TEXT DEFAULT (datetime('now')),
		updated_at TEXT DEFAULT (datetime('now'))
	);
	CREATE TABLE IF NOT EXISTS collection_books (
		collection_id INTEGER,
		book_id INTEGER,
		added_at TEXT DEFAULT (datetime('now')),
		PRIMARY KEY (collection_id, book_id),
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
		FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	return db
}

func TestGetBooks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-03-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/api/v1/books", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var books []models.Book
	err = json.NewDecoder(rr.Body).Decode(&books)
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)
	assert.Equal(t, "Test Author", books[0].Author)
	assert.Equal(t, "2022-03-01", books[0].PublishedDate)
	assert.Equal(t, "1st", books[0].Edition)
	assert.Equal(t, "A test book description", books[0].Description)
	assert.Equal(t, "Test Genre", books[0].Genre)
}

func TestGetBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-03-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/api/v1/books/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var book models.Book
	t.Log(rr.Body.String())
	err = json.NewDecoder(rr.Body).Decode(&book)
	assert.NoError(t, err)
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, "Test Author", book.Author)
	assert.Equal(t, "2022-03-01", book.PublishedDate)
	assert.Equal(t, "1st", book.Edition)
	assert.Equal(t, "A test book description", book.Description)
	assert.Equal(t, "Test Genre", book.Genre)
}

func TestCreateBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	book := models.Book{
		Title:         "New Book",
		Author:        "New Author",
		PublishedDate: "2022-01-01",
		Edition:       "1st",
		Description:   "A new book description",
		Genre:         "New Genre",
	}
	bookJSON, _ := json.Marshal(book)

	req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(bookJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdBook models.Book
	err = json.NewDecoder(rr.Body).Decode(&createdBook)
	t.Log(rr.Body.String())

	assert.NoError(t, err)
	assert.Equal(t, "New Book", createdBook.Title)
	assert.Equal(t, "New Author", createdBook.Author)
	assert.Equal(t, "2022-01-01", createdBook.PublishedDate)
	assert.Equal(t, "1st", createdBook.Edition)
	assert.Equal(t, "A new book description", createdBook.Description)
	assert.Equal(t, "New Genre", createdBook.Genre)

	// Verify insertion count
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books WHERE title = ?", "New Book").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify values
	var title, author, edition, description, genre string
	var publishedDate, createdAt, updatedAt string
	err = db.QueryRow("SELECT title, author, published_date, edition, description, genre, created_at, updated_at FROM books WHERE id = 1", "New Book").Scan(&title, &author, &publishedDate, &edition, &description, &genre, &createdAt, &updatedAt)
	assert.NoError(t, err)
	assert.Equal(t, "New Book", title)
	assert.Equal(t, "New Author", author)
	assert.Equal(t, "2022-01-01", publishedDate)
	assert.Equal(t, "1st", edition)
	assert.Equal(t, "A new book description", description)
	assert.Equal(t, "New Genre", genre)
	assert.NotEmpty(t, createdAt)
	assert.NotEmpty(t, updatedAt)
	assert.NotEqual(t, "0001-01-01T00:00:00Z", createdAt)
	assert.NotEqual(t, "0001-01-01T00:00:00Z", updatedAt)
}

func TestUpdateBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	// Ensure updated_at is different
	time.Sleep(1 * time.Second)

	book := models.Book{
		Title:         "Updated Book",
		Author:        "Updated Author",
		PublishedDate: "2022-01-01",
		Edition:       "2nd",
		Description:   "An updated book description",
		Genre:         "Updated Genre",
	}
	bookJSON, _ := json.Marshal(book)

	req, err := http.NewRequest("PUT", "/api/v1/books/1", bytes.NewBuffer(bookJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var updatedBook models.Book
	err = json.NewDecoder(rr.Body).Decode(&updatedBook)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Book", updatedBook.Title)
	assert.Equal(t, "Updated Author", updatedBook.Author)
	assert.Equal(t, "2022-01-01", updatedBook.PublishedDate)
	assert.Equal(t, "2nd", updatedBook.Edition)
	assert.Equal(t, "An updated book description", updatedBook.Description)
	assert.Equal(t, "Updated Genre", updatedBook.Genre)

	// Verify data
	var title, author, edition, description, genre, publishedDate, createdAt, updatedAt string
	err = db.QueryRow("SELECT title, author, published_date, edition, description, genre, created_at, updated_at FROM books WHERE id = 1", "Updated Book").Scan(&title, &author, &publishedDate, &edition, &description, &genre, &createdAt, &updatedAt)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Book", title)
	assert.Equal(t, "Updated Author", author)
	assert.Equal(t, "2022-01-01", publishedDate)
	assert.Equal(t, "2nd", edition)
	assert.Equal(t, "An updated book description", description)
	assert.Equal(t, "Updated Genre", genre)
	assert.NotEmpty(t, createdAt)
	assert.NotEmpty(t, updatedAt)
	assert.NotEqual(t, "0001-01-01T00:00:00Z", createdAt)
	assert.NotEqual(t, "0001-01-01T00:00:00Z", updatedAt)
	assert.NotEqual(t, createdAt, updatedAt)
}

func TestDeleteBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/api/v1/books/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books WHERE id = ?", 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestGetCollections(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/api/v1/collections", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var collections []models.Collection
	err = json.NewDecoder(rr.Body).Decode(&collections)
	assert.NoError(t, err)
	assert.Len(t, collections, 1)
	assert.Equal(t, "Test Collection", collections[0].Name)
}

func TestGetCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO collection_books (collection_id, book_id) VALUES (?, ?)", 1, 1)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "/api/v1/collections/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var collection models.Collection
	err = json.NewDecoder(rr.Body).Decode(&collection)
	assert.NoError(t, err)
	assert.Equal(t, "Test Collection", collection.Name)
	assert.Len(t, collection.Books, 1)
	assert.Equal(t, "Test Book", collection.Books[0].Title)
	assert.Equal(t, "Test Author", collection.Books[0].Author)
	assert.Equal(t, "2022-01-01", collection.Books[0].PublishedDate)
	assert.Equal(t, "1st", collection.Books[0].Edition)
	assert.Equal(t, "A test book description", collection.Books[0].Description)
	assert.Equal(t, "Test Genre", collection.Books[0].Genre)
}

func TestCreateCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	collection := models.Collection{Name: "New Collection"}
	collectionJSON, _ := json.Marshal(collection)

	req, err := http.NewRequest("POST", "/api/v1/collections", bytes.NewBuffer(collectionJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var createdCollection models.Collection
	err = json.NewDecoder(rr.Body).Decode(&createdCollection)
	assert.NoError(t, err)
	assert.Equal(t, "New Collection", createdCollection.Name)

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collections WHERE name = ?", "New Collection").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	var name, createdAt, updatedAt string
	err = db.QueryRow("SELECT name, created_at, updated_at FROM collections WHERE id = 1").Scan(&name, &createdAt, &updatedAt)
	assert.NoError(t, err)
	assert.Equal(t, "New Collection", name)
	assert.NotEmpty(t, createdAt)
	assert.NotEmpty(t, updatedAt)
	assert.NotEqual(t, "0001-01-01T00:00:00Z", createdAt)
	assert.NotEqual(t, "0001-01-01T00:00:00Z", updatedAt)
}

func TestUpdateCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)

	collection := models.Collection{Name: "Updated Collection"}
	collectionJSON, _ := json.Marshal(collection)

	req, err := http.NewRequest("PUT", "/api/v1/collections/1", bytes.NewBuffer(collectionJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var updatedCollection models.Collection
	err = json.NewDecoder(rr.Body).Decode(&updatedCollection)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Collection", updatedCollection.Name)

	// Verify update
	var updatedName string
	err = db.QueryRow("SELECT name FROM collections WHERE id = ?", 1).Scan(&updatedName)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Collection", updatedName)
}

func TestDeleteCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/api/v1/collections/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collections WHERE id = ?", 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestAddBookToCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/collections/1/books/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collection_books WHERE collection_id = ? AND book_id = ?", 1, 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestRemoveBookFromCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO collection_books (collection_id, book_id) VALUES (?, ?)", 1, 1)
	assert.NoError(t, err)

	req, err := http.NewRequest("DELETE", "/api/v1/collections/1/books/1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupTestRouter(db)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collection_books WHERE collection_id = ? AND book_id = ?", 1, 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
