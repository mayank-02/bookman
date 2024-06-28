package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mayank-02/bookman/internal/models"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Manage book collections",
}

var collectionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new collection",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		collection := models.Collection{Name: name}

		createdCollection, err := bookman.CreateCollection(collection)
		handleErr(err)
		printCollectionDetails(createdCollection)
	},
}

var collectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all collections",
	Run: func(cmd *cobra.Command, args []string) {
		collections, err := bookman.GetCollections()
		handleErr(err)
		printCollectionsTable(collections)
	},
}

var collectionGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details of a specific collection",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		collectionID, err := strconv.Atoi(id)
		handleErr(err)

		collection, err := bookman.GetCollection(collectionID)
		handleErr(err)
		printCollectionDetails(collection)
	},
}

var collectionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a collection",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		collectionID, err := strconv.Atoi(id)
		handleErr(err)
		name, _ := cmd.Flags().GetString("name")

		collection := models.Collection{
			ID:   collectionID,
			Name: name,
		}

		err = bookman.UpdateCollection(collection)
		handleErr(err)
		fmt.Println("Collection updated successfully")
	},
}

var collectionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a collection",
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetString("id")
		collectionID, err := strconv.Atoi(id)
		handleErr(err)

		err = bookman.DeleteCollection(collectionID)
		handleErr(err)
		fmt.Println("Collection deleted successfully")
	},
}

var collectionAddBookCmd = &cobra.Command{
	Use:   "add-book",
	Short: "Add a book to a collection",
	Run: func(cmd *cobra.Command, args []string) {
		collectionID, _ := cmd.Flags().GetString("collection-id")
		bookID, _ := cmd.Flags().GetString("book-id")
		colID, err := strconv.Atoi(collectionID)
		handleErr(err)
		bkID, err := strconv.Atoi(bookID)
		handleErr(err)

		err = bookman.AddBookToCollection(colID, bkID)
		handleErr(err)
		fmt.Println("Book added to collection successfully")
	},
}

var collectionRemoveBookCmd = &cobra.Command{
	Use:   "remove-book",
	Short: "Remove a book from a collection",
	Run: func(cmd *cobra.Command, args []string) {
		collectionID, _ := cmd.Flags().GetString("collection-id")
		bookID, _ := cmd.Flags().GetString("book-id")
		colID, err := strconv.Atoi(collectionID)
		handleErr(err)
		bkID, err := strconv.Atoi(bookID)
		handleErr(err)

		err = bookman.RemoveBookFromCollection(colID, bkID)
		handleErr(err)
		fmt.Println("Book removed from collection successfully")
	},
}

func printCollectionsTable(collections []models.Collection) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})

	for _, collection := range collections {
		table.Append([]string{
			strconv.Itoa(collection.ID),
			collection.Name,
		})
	}

	table.Render()
}

func printCollectionDetails(collection models.Collection) {
	fmt.Println("Collection:")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name"})

	table.Append([]string{strconv.Itoa(collection.ID), collection.Name})

	table.Render()

	// Print books in the collection
	if len(collection.Books) > 0 {
		fmt.Println("\nBooks in the collection:")
		printBooksTable(collection.Books)
	}
}

func init() {
	collectionCreateCmd.Flags().String("name", "", "Name of the collection")

	collectionGetCmd.Flags().String("id", "", "ID of the collection")
	collectionUpdateCmd.Flags().String("id", "", "ID of the collection")
	collectionUpdateCmd.Flags().String("name", "", "Name of the collection")
	collectionDeleteCmd.Flags().String("id", "", "ID of the collection")

	collectionAddBookCmd.Flags().String("collection-id", "", "ID of the collection")
	collectionAddBookCmd.Flags().String("book-id", "", "ID of the book")
	collectionRemoveBookCmd.Flags().String("collection-id", "", "ID of the collection")
	collectionRemoveBookCmd.Flags().String("book-id", "", "ID of the book")

	collectionCmd.AddCommand(collectionCreateCmd)
	collectionCmd.AddCommand(collectionListCmd)
	collectionCmd.AddCommand(collectionGetCmd)
	collectionCmd.AddCommand(collectionUpdateCmd)
	collectionCmd.AddCommand(collectionDeleteCmd)
	collectionCmd.AddCommand(collectionAddBookCmd)
	collectionCmd.AddCommand(collectionRemoveBookCmd)
}
