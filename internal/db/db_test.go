package db

import (
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mayank-02/bookman/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := InitDB(":memory:")
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

func TestDB_GetBooks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	// Retrieve books without filters
	books, err := db.GetBooks("", "", "", "")
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)
	assert.Equal(t, "Test Author", books[0].Author)
	assert.Equal(t, "2022-01-01", books[0].PublishedDate)
	assert.Equal(t, "1st", books[0].Edition)
	assert.Equal(t, "A test book description", books[0].Description)
	assert.Equal(t, "Test Genre", books[0].Genre)

	// Retrieve books with author filter
	books, err = db.GetBooks("Test Author", "", "", "")
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)

	// Retrieve books with genre filter
	books, err = db.GetBooks("", "Test Genre", "", "")
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)

	// Retrieve books with date range filter
	books, err = db.GetBooks("", "", "2022-01-01", "2022-12-31")
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, "Test Book", books[0].Title)

	// Retrieve books with non-matching filters
	books, err = db.GetBooks("Non-existent Author", "", "", "")
	assert.NoError(t, err)
	assert.Len(t, books, 0)
}

func TestDB_GetBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	// Retrieve book
	book, err := db.GetBook(1)
	assert.NoError(t, err)
	assert.Equal(t, "Test Book", book.Title)
	assert.Equal(t, "Test Author", book.Author)
	assert.Equal(t, "2022-01-01", book.PublishedDate)
	assert.Equal(t, "1st", book.Edition)
	assert.Equal(t, "A test book description", book.Description)
	assert.Equal(t, "Test Genre", book.Genre)
}

func TestDB_CreateBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create book
	book := models.Book{
		Title:         "New Book",
		Author:        "New Author",
		PublishedDate: "2022-01-01",
		Edition:       "1st",
		Description:   "A new book description",
		Genre:         "New Genre",
	}
	id, err := db.CreateBook(book)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books WHERE title = ?", "New Book").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify timestamps
	var createdAt, updatedAt string
	err = db.QueryRow("SELECT created_at, updated_at FROM books WHERE id = ?", 1).Scan(&createdAt, &updatedAt)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdAt)
	assert.NotEmpty(t, updatedAt)
	assert.Equal(t, createdAt, updatedAt)

	// Verify data
	var title, author, publishedDate, edition, description, genre string
	err = db.QueryRow("SELECT title, author, published_date, edition, description, genre FROM books WHERE id = ?", 1).Scan(&title, &author, &publishedDate, &edition, &description, &genre)
	assert.NoError(t, err)
	assert.Equal(t, "New Book", title)
	assert.Equal(t, "New Author", author)
	assert.Equal(t, "2022-01-01", publishedDate)
	assert.Equal(t, "1st", edition)
	assert.Equal(t, "A new book description", description)
	assert.Equal(t, "New Genre", genre)
}

func TestDB_UpdateBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	// Sleep for 1 second to ensure updated_at timestamp is different
	time.Sleep(1 * time.Second)

	// Update book
	book := models.Book{
		ID:            1,
		Title:         "Updated Book",
		Author:        "Updated Author",
		PublishedDate: "2022-01-01",
		Edition:       "2nd",
		Description:   "An updated book description",
		Genre:         "Updated Genre",
	}
	err = db.UpdateBook(book)
	assert.NoError(t, err)

	// Verify update
	var title, author, publishedDate, edition, description, genre string
	err = db.QueryRow("SELECT title, author, published_date, edition, description, genre FROM books WHERE id = ?", 1).Scan(&title, &author, &publishedDate, &edition, &description, &genre)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Book", title)
	assert.Equal(t, "Updated Author", author)
	assert.Equal(t, "2022-01-01", publishedDate)
	assert.Equal(t, "2nd", edition)
	assert.Equal(t, "An updated book description", description)
	assert.Equal(t, "Updated Genre", genre)

	// Verify timestamps
	var createdAt, updatedAt string
	err = db.QueryRow("SELECT created_at, updated_at FROM books WHERE id = ?", 1).Scan(&createdAt, &updatedAt)
	assert.NoError(t, err)
	assert.NotEmpty(t, createdAt)
	assert.NotEmpty(t, updatedAt)
	assert.NotEqual(t, createdAt, updatedAt)
}

func TestDB_DeleteBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	// Delete book
	err = db.DeleteBook(1)
	assert.NoError(t, err)

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books WHERE id = ?", 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDB_GetCollections(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)

	// Retrieve collections
	collections, err := db.GetCollections()
	assert.NoError(t, err)
	assert.Len(t, collections, 1)
	assert.Equal(t, "Test Collection", collections[0].Name)
}

func TestDB_GetCollection(t *testing.T) {
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

	// Retrieve collection
	collection, err := db.GetCollection(1)
	assert.NoError(t, err)
	assert.Equal(t, "Test Collection", collection.Name)
	assert.Len(t, collection.Books, 1)
	assert.Equal(t, "Test Book", collection.Books[0].Title)
}

func TestDB_CreateCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create collection
	collection := models.Collection{Name: "New Collection"}
	id, err := db.CreateCollection(collection)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collections WHERE name = ?", "New Collection").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestDB_UpdateCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)

	// Update collection
	collection := models.Collection{
		ID:   1,
		Name: "Updated Collection",
	}
	err = db.UpdateCollection(collection)
	assert.NoError(t, err)

	// Verify update
	var updatedName string
	err = db.QueryRow("SELECT name FROM collections WHERE id = ?", 1).Scan(&updatedName)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Collection", updatedName)
}

func TestDB_DeleteCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)

	// Delete collection
	err = db.DeleteCollection(1)
	assert.NoError(t, err)

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collections WHERE id = ?", 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestDB_AddBookToCollection(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Insert test data
	_, err := db.Exec("INSERT INTO collections (name) VALUES (?)", "Test Collection")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		"Test Book", "Test Author", "2022-01-01", "1st", "A test book description", "Test Genre")
	assert.NoError(t, err)

	// Add book to collection
	err = db.AddBookToCollection(1, 1)
	assert.NoError(t, err)

	// Verify insertion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collection_books WHERE collection_id = ? AND book_id = ?", 1, 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestDB_RemoveBookFromCollection(t *testing.T) {
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

	// Remove book from collection
	err = db.RemoveBookFromCollection(1, 1)
	assert.NoError(t, err)

	// Verify deletion
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM collection_books WHERE collection_id = ? AND book_id = ?", 1, 1).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
