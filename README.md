# Bookman

Bookman is a simple book management software designed to provide users with an efficient way to manage their book collections. The software offers a robust REST API, utilizes a relational database (SQLite) for data storage, and includes a command-line interface (CLI) for easy interaction.

## Features

As a user, you can:
- Add and manage books into the system, including some basic information about those books (title, author, published date, edition, description, genre, ...)
- Create and manage collections of books
- Easily list all books, all collections, and filter book lists by author, genre, or a range of publication dates

## Setup

Pre-requisites: Go 1.22 or higher, SQLite and Git

Steps to run the application:
```bash
# Clone the Repository
git clone https://github.com/mayank-02/bookman.git

# Change the working directory
cd bookman

# Install dependencies
go mod download

# Create the SQLite database
sqlite3 bookman.db < sql/schema.sql

# Ensure tests pass
go test ./...

# Build and run the server
go build -o bin/server cmd/server/main.go
./bin/server # The server will start on `http://localhost:8080`.

# Build and run the CLI
go build -o bin/bookman cmd/cli/*
./bin/bookman 
```

## Usage

Below is a comprehensive guide to using the CLI, detailing all available commands, their options, and examples of user interactions.

```bash
bookman - Book Management CLI

Usage:
  bookman [command]

Available Commands:
  book        Manage books
  collection  Manage book collections
  help        Help about any command
  version     Print the version number of bookman

Book Commands:
  book add         Add a new book
  book delete      Delete a book
  book get         Get details of a specific book
  book list        List all books
  book update      Update a book's information

Collection Commands:
  collection add-book       Add a book to a collection
  collection create         Create a new collection
  collection delete         Delete a collection
  collection get            Get details of a specific collection
  collection list           List all collections
  collection remove-book    Remove a book from a collection
  collection update         Update a collection

Flags:
  -h, --help      Help for bookman

Use "bookman [command] --help" for more information about a command.
```

Book-related commands:
```bash
# Adding a book
$ bookman book add --title "The Go Programming Language" --author "Alan A. A. Donovan, Brian W. Kernighan" --published "2015-10-26" --genre "Programming" --description "An authoritative resource for Go programming language"

# Getting details of a book
$ bookman book get --id 1

# Listing books
$ bookman book list
$ bookman book list --author "Alan A. A. Donovan, Brian W. Kernighan"
$ bookman book list --genre "Programming"
$ bookman book list --from "2010-01-01" --to "2020-12-31"

# Updating a book
$ bookman book update --id 1 --title "The Go Programming Language (2nd Edition)" --author "Mayank Jain" --published "2024-06-28"

# Deleting a book
$ bookman book delete --id 1
```

Collection-related commands:
```bash
# Creating a collection
$ bookman collection create --name "My Favorite Programming Books"

# Adding a book to a collection
$ bookman collection add-book --collection-id 1 --book-id 1

# Listing collections
$ bookman collection list

# Getting details of a collection
$ bookman collection get --id 1

# Removing a book from a collection
$ bookman collection remove-book --collection-id 1 --book-id 1
```
## REST API

### Models

#### Book

```json
{
  "id": 1,
  "title": "string",
  "author": "string",
  "published_date": "YYYY-MM-DD",
  "edition": "string",
  "description": "string",
  "genre": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

#### Collection

```json
{
  "id": 1,
  "name": "string",
  "books": [
    {
      "id": 1,
      "title": "string",
      "author": "string",
      "published_date": "YYYY-MM-DD",
      "edition": "string",
      "description": "string",
      "genre": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  ]
}
```

### Books API

| Method | Endpoint           | Description              | Request Body                                                                                                                                 | Query Parameters                                                            | Response Code | Response Body |
| ------ | ------------------ | ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------- | ------------- | ------------- |
| GET    | /api/v1/books      | Retrieve all books       | N/A                                                                                                                                          | `author` (optional), `genre` (optional), `from` (optional), `to` (optional) | 200           | List\<Book\>  |
| POST   | /api/v1/books      | Create a new book        | `{ "title": "string", "author": "string", "published_date": "YYYY-MM-DD", "edition": "string", "description": "string", "genre": "string" }` | N/A                                                                         | 201           | Book          |
| GET    | /api/v1/books/{id} | Retrieve a specific book | N/A                                                                                                                                          | N/A                                                                         | 200           | Book          |
| PUT    | /api/v1/books/{id} | Update a specific book   | `{ "title": "string", "author": "string", "published_date": "YYYY-MM-DD", "edition": "string", "description": "string", "genre": "string" }` | N/A                                                                         | 200           | Book          |
| DELETE | /api/v1/books/{id} | Delete a specific book   | N/A                                                                                                                                          | N/A                                                                         | 204           | N/A           |

### Collections API

| Method | Endpoint                                | Description                              | Request Body           | Response Code | Response Body                                |
| ------ | --------------------------------------- | ---------------------------------------- | ---------------------- | ------------- | -------------------------------------------- |
| GET    | /api/v1/collections                     | Retrieve all collections                 | N/A                    | 200           | List\<Collection\>                           |
| POST   | /api/v1/collections                     | Create a new collection                  | `{ "name": "string" }` | 201           | `{ "id": 1, "name": "string", "books": [] }` |
| GET    | /api/v1/collections/{id}                | Retrieve a specific collection           | N/A                    | 200           | Collection                                   |
| PUT    | /api/v1/collections/{id}                | Update a specific collection             | `{ "name": "string" }` | 200           | Collection                                   |
| DELETE | /api/v1/collections/{id}                | Delete a specific collection             | N/A                    | 204           | N/A                                          |
| POST   | /api/v1/collections/{id}/books/{bookId} | Add a book to a specific collection      | N/A                    | 204           | N/A                                          |
| DELETE | /api/v1/collections/{id}/books/{bookId} | Remove a book from a specific collection | N/A                    | 204           | N/A                                          |


### Error Handling

All endpoints return appropriate HTTP status codes and error messages in the following cases:

- 400 Bad Request: Invalid input or missing required fields
- 404 Not Found: Resource not found
- 409 Conflict: Conflict in the request, e.g., adding a book that is already in the collection
- 500 Internal Server Error: Internal server error


## Database Schema

This should include all the tables you expect to use to store this data, for each table, all its
expected columns, column types, relations, constraints and any index you’d like to add.
You may write this directly in SQL or use tables or text as an alternative way to describe what
you’re going for.

```text
+-------------------+             +--------------------------+             +-------------------+
|    books          |             |     collection_books     |             |    collections    |
+-------------------+             +--------------------------+             +-------------------+
| - id (PK)         |<-----┐      | - collection_id (FK, PK) |<----------->| - id (PK)         |
| - title           |      └----->| - book_id (FK, PK)       |             | - name            |
| - author          |             | - added_at               |             +-------------------+
| - published_date  |             +--------------------------+
| - edition         |
| - description     |
| - genre           |
| - created_at      |
| - updated_at      |
+-------------------+

Indexes: On author, genre, published_date in books table and on collection_id, book_id in collection_books table.
```

## Directory Structure

```
bookman
├── LICENSE
├── README.md                     # Overview and usage instructions
├── cmd
│   ├── cli                       # CLI related commands
│   │   ├── book.go
│   │   ├── collection.go
│   │   └── main.go               # Entry point for the CLI application
│   └── server                    # Server related commands
│       └── main.go               # Entry point for the server application
├── go.mod
├── go.sum
├── internal
│   ├── api
│   │   ├── handlers.go           # API endpoint handlers
│   │   └── handlers_test.go      # Tests for API handlers
│   ├── db
│   │   ├── db.go                 # Database initialization and operations
│   │   └── db_test.go            # Tests for database operations
│   └── models                    # Data models
│       ├── book.go
│       └── collection.go
├── pkg
│   └── client                    # Client package for interacting with the server
│       └── client.go
└── sql
    └── schema.sql                # SQL schema for setting up the database
```



## Possible Enhancements

1. Secure the API with user authentication and authorization
2. Add pagination support for book listings
3. Implement advanced search capabilities and multi-criteria filtering
4. Integrate with external APIs like Google Books for additional information
5. Optimize database queries and indexing for large datasets

## Resources

1. Golang Documentation: [golang.org/doc](https://golang.org/doc/)
2. SQLite Documentation: [sqlite.org/docs](https://www.sqlite.org/docs.html)
3. Gorilla Mux: [github.com/gorilla/mux](https://github.com/gorilla/mux)
4. Cobra CLI: [github.com/spf13/cobra](https://github.com/spf13/cobra)
5. Testing in Go: [golang.org/pkg/testing](https://golang.org/pkg/testing/), [golang.org/pkg/net/http/httptest](https://golang.org/pkg/net/http/httptest/) and [stretchr/testify](https://github.com/stretchr/testify)

## Author

[Mayank Jain](https://jainmayank.me/)