package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mayank-02/bookman/internal/models"
)

type DB struct {
	*sql.DB
}

func InitDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetBook(id int) (models.Book, error) {
	var b models.Book
	var createdAt, updatedAt string
	err := db.QueryRow("SELECT id, title, author, published_date, edition, description, genre, created_at, updated_at FROM books WHERE id = ?", id).
		Scan(&b.ID, &b.Title, &b.Author, &b.PublishedDate, &b.Edition, &b.Description, &b.Genre, &createdAt, &updatedAt)
	if err != nil {
		return models.Book{}, err
	}

	b.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return models.Book{}, err
	}
	b.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
	if err != nil {
		return models.Book{}, err
	}
	return b, nil
}

func (db *DB) GetBooks(author, genre, from, to string) ([]models.Book, error) {
	query := "SELECT id, title, author, published_date, edition, description, genre, created_at, updated_at FROM books WHERE 1=1"
	args := []interface{}{}

	if author != "" {
		query += " AND author = ?"
		args = append(args, author)
	}
	if genre != "" {
		query += " AND genre = ?"
		args = append(args, genre)
	}
	if from != "" {
		query += " AND published_date >= ?"
		args = append(args, from)
	}
	if to != "" {
		query += " AND published_date <= ?"
		args = append(args, to)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		var createdAt, updatedAt string
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedDate, &b.Edition, &b.Description, &b.Genre, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		b.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, err
		}
		b.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (db *DB) CreateBook(b models.Book) (int, error) {
	result, err := db.Exec("INSERT INTO books (title, author, published_date, edition, description, genre) VALUES (?, ?, ?, ?, ?, ?)",
		b.Title, b.Author, b.PublishedDate, b.Edition, b.Description, b.Genre)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *DB) UpdateBook(b models.Book) error {
	_, err := db.Exec("UPDATE books SET title = ?, author = ?, published_date = ?, edition = ?, description = ?, genre = ?, updated_at = datetime('now') WHERE id = ?",
		b.Title, b.Author, b.PublishedDate, b.Edition, b.Description, b.Genre, b.ID)
	return err
}

func (db *DB) DeleteBook(id int) error {
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	return err
}

func (db *DB) GetCollections() ([]models.Collection, error) {
	rows, err := db.Query("SELECT id, name FROM collections")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []models.Collection
	for rows.Next() {
		var c models.Collection
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}
	return collections, nil
}

func (db *DB) GetCollection(id int) (models.Collection, error) {
	var c models.Collection
	err := db.QueryRow("SELECT id, name FROM collections WHERE id = ?", id).Scan(&c.ID, &c.Name)
	if err != nil {
		return models.Collection{}, err
	}

	rows, err := db.Query(`
		SELECT b.id, b.title, b.author, b.published_date, b.edition, b.description, b.genre, b.created_at, b.updated_at
		FROM books b
		JOIN collection_books cb ON b.id = cb.book_id
		WHERE cb.collection_id = ?`, id)
	if err != nil {
		return models.Collection{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var b models.Book
		var createdAt, updatedAt string
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedDate, &b.Edition, &b.Description, &b.Genre, &createdAt, &updatedAt)
		if err != nil {
			return models.Collection{}, err
		}
		b.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return models.Collection{}, err
		}
		b.UpdatedAt, err = time.Parse("2006-01-02 15:04:05", updatedAt)
		if err != nil {
			return models.Collection{}, err
		}
		c.Books = append(c.Books, b)
	}

	return c, nil
}

func (db *DB) CreateCollection(c models.Collection) (int, error) {
	result, err := db.Exec("INSERT INTO collections (name) VALUES (?)", c.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (db *DB) UpdateCollection(c models.Collection) error {
	_, err := db.Exec("UPDATE collections SET name = ? WHERE id = ?", c.Name, c.ID)
	return err
}

func (db *DB) DeleteCollection(id int) error {
	_, err := db.Exec("DELETE FROM collections WHERE id = ?", id)
	return err
}

func (db *DB) AddBookToCollection(collectionID, bookID int) error {
	_, err := db.Exec("INSERT INTO collection_books (collection_id, book_id) VALUES (?, ?)", collectionID, bookID)
	return err
}

func (db *DB) RemoveBookFromCollection(collectionID, bookID int) error {
	_, err := db.Exec("DELETE FROM collection_books WHERE collection_id = ? AND book_id = ?", collectionID, bookID)
	return err
}

func (db *DB) IsBookInCollection(collectionID, bookID int) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM collection_books WHERE collection_id = ? AND book_id = ?", collectionID, bookID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
