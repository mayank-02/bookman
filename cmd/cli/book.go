package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mayank-02/bookman/internal/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Manage books",
}

func parseBookFlags(cmd *cobra.Command) (models.Book, error) {
	title, _ := cmd.Flags().GetString("title")
	author, _ := cmd.Flags().GetString("author")
	publishedDate, _ := cmd.Flags().GetString("published")
	edition, _ := cmd.Flags().GetString("edition")
	description, _ := cmd.Flags().GetString("description")
	genre, _ := cmd.Flags().GetString("genre")

	if title == "" || author == "" || publishedDate == "" {
		return models.Book{}, fmt.Errorf("title, author, and published are mandatory fields")
	}

	return models.Book{
		Title:         title,
		Author:        author,
		PublishedDate: publishedDate,
		Edition:       edition,
		Description:   description,
		Genre:         genre,
	}, nil
}

func printBooksTable(books []models.Book) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Author", "Published Date", "Edition", "Description", "Genre"})

	for _, book := range books {
		table.Append([]string{
			strconv.Itoa(book.ID),
			book.Title,
			book.Author,
			book.PublishedDate,
			book.Edition,
			book.Description,
			book.Genre,
		})
	}

	table.Render()
}

var bookAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new book",
	Run: func(cmd *cobra.Command, args []string) {
		book, err := parseBookFlags(cmd)
		handleErr(err)

		_, err = bookman.CreateBook(book)
		handleErr(err)
		fmt.Println("Book added successfully")
	},
}

var bookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all books",
	Run: func(cmd *cobra.Command, args []string) {
		author, _ := cmd.Flags().GetString("author")
		genre, _ := cmd.Flags().GetString("genre")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		books, err := bookman.GetBooks(author, genre, from, to)
		handleErr(err)
		printBooksTable(books)

	},
}

var bookGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific book",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		bookID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid ID format:", err)
			return
		}

		book, err := bookman.GetBook(bookID)
		handleErr(err)
		printBooksTable([]models.Book{book})
	},
}

var bookUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a book's information",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		bookID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid ID format:", err)
			return
		}

		book, err := parseBookFlags(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
		book.ID = bookID

		err = bookman.UpdateBook(book)
		handleErr(err)
		fmt.Println("Book updated successfully")
	},
}

var bookDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a book",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		bookID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid ID format:", err)
			return
		}

		err = bookman.DeleteBook(bookID)
		handleErr(err)
		fmt.Println("Book deleted successfully")
	},
}

func init() {
	bookAddCmd.Flags().String("title", "", "Title of the book")
	bookAddCmd.Flags().String("author", "", "Author of the book")
	bookAddCmd.Flags().String("published", "", "Published date of the book")
	bookAddCmd.Flags().String("edition", "", "Edition of the book")
	bookAddCmd.Flags().String("description", "", "Description of the book")
	bookAddCmd.Flags().String("genre", "", "Genre of the book")

	bookListCmd.Flags().String("author", "", "Filter books by author")
	bookListCmd.Flags().String("genre", "", "Filter books by genre")
	bookListCmd.Flags().String("from", "", "Filter books published after this date")
	bookListCmd.Flags().String("to", "", "Filter books published before this date")

	bookGetCmd.Flags().String("id", "", "ID of the book")

	bookUpdateCmd.Flags().String("id", "", "ID of the book")
	bookUpdateCmd.Flags().String("title", "", "Title of the book")
	bookUpdateCmd.Flags().String("author", "", "Author of the book")
	bookUpdateCmd.Flags().String("published", "", "Published date of the book")
	bookUpdateCmd.Flags().String("edition", "", "Edition of the book")
	bookUpdateCmd.Flags().String("description", "", "Description of the book")
	bookUpdateCmd.Flags().String("genre", "", "Genre of the book")

	bookDeleteCmd.Flags().String("id", "", "ID of the book")

	bookCmd.AddCommand(bookAddCmd, bookListCmd, bookGetCmd, bookUpdateCmd, bookDeleteCmd)
}
