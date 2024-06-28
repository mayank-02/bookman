package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mayank-02/bookman/internal/models"
)

type Client struct {
	BaseURL    string
	HttpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetBooks(author, genre, from, to string) ([]models.Book, error) {
	// Build the query parameters
	queryParams := url.Values{}
	if author != "" {
		queryParams.Add("author", author)
	}
	if genre != "" {
		queryParams.Add("genre", genre)
	}
	if from != "" {
		queryParams.Add("from", from)
	}
	if to != "" {
		queryParams.Add("to", to)
	}

	// Build the complete URL with query parameters
	apiURL := c.BaseURL + "/api/v1/books"
	if len(queryParams) > 0 {
		apiURL += "?" + queryParams.Encode()
	}

	// Make the HTTP GET request
	resp, err := c.HttpClient.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response into the books slice
	var books []models.Book
	err = json.NewDecoder(resp.Body).Decode(&books)
	return books, err
}

func (c *Client) GetBook(id int) (models.Book, error) {
	resp, err := c.HttpClient.Get(fmt.Sprintf("%s/api/v1/books/%d", c.BaseURL, id))
	if err != nil {
		return models.Book{}, err
	}
	defer resp.Body.Close()

	var book models.Book
	err = json.NewDecoder(resp.Body).Decode(&book)
	return book, err
}

func (c *Client) CreateBook(book models.Book) (models.Book, error) {
	bookJSON, _ := json.Marshal(book)
	resp, err := c.HttpClient.Post(c.BaseURL+"/api/v1/books", "application/json", bytes.NewBuffer(bookJSON))
	if err != nil {
		return models.Book{}, err
	}
	defer resp.Body.Close()

	var createdBook models.Book
	err = json.NewDecoder(resp.Body).Decode(&createdBook)
	return createdBook, err
}

func (c *Client) UpdateBook(book models.Book) error {
	bookJSON, _ := json.Marshal(book)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/books/%d", c.BaseURL, book.ID), bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update book: %s", resp.Status)
	}
	return nil
}

func (c *Client) DeleteBook(id int) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/books/%d", c.BaseURL, id), nil)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete book: %s", resp.Status)
	}
	return nil
}

func (c *Client) GetCollections() ([]models.Collection, error) {
	resp, err := c.HttpClient.Get(c.BaseURL + "/api/v1/collections")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var collections []models.Collection
	err = json.NewDecoder(resp.Body).Decode(&collections)
	return collections, err
}

func (c *Client) GetCollection(id int) (models.Collection, error) {
	resp, err := c.HttpClient.Get(fmt.Sprintf("%s/api/v1/collections/%d", c.BaseURL, id))
	if err != nil {
		return models.Collection{}, err
	}
	defer resp.Body.Close()

	var collection models.Collection
	err = json.NewDecoder(resp.Body).Decode(&collection)
	return collection, err
}

func (c *Client) CreateCollection(collection models.Collection) (models.Collection, error) {
	collectionJSON, _ := json.Marshal(collection)
	resp, err := c.HttpClient.Post(c.BaseURL+"/api/v1/collections", "application/json", bytes.NewBuffer(collectionJSON))
	if err != nil {
		return models.Collection{}, err
	}
	defer resp.Body.Close()

	var createdCollection models.Collection
	err = json.NewDecoder(resp.Body).Decode(&createdCollection)
	return createdCollection, err
}

func (c *Client) UpdateCollection(collection models.Collection) error {
	collectionJSON, _ := json.Marshal(collection)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/collections/%d", c.BaseURL, collection.ID), bytes.NewBuffer(collectionJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update collection: %s", resp.Status)
	}
	return nil
}

func (c *Client) DeleteCollection(id int) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/collections/%d", c.BaseURL, id), nil)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete collection: %s", resp.Status)
	}
	return nil
}

func (c *Client) AddBookToCollection(collectionID, bookID int) error {
	url := fmt.Sprintf("%s/api/v1/collections/%d/books/%d", c.BaseURL, collectionID, bookID)
	resp, err := c.HttpClient.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add book to collection: %s", string(body))
	}
	return nil
}

func (c *Client) RemoveBookFromCollection(collectionID, bookID int) error {
	url := fmt.Sprintf("%s/api/v1/collections/%d/books/%d", c.BaseURL, collectionID, bookID)
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove book from collection: %s", string(body))
	}
	return nil
}
