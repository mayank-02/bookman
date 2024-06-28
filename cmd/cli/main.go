package main

import (
	"fmt"
	"os"

	"github.com/mayank-02/bookman/pkg/client"
	"github.com/spf13/cobra"
)

var (
	bookman = client.New("http://localhost:8080")
)

const Version = "v1.0.0"

func handleErr(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func main() {
	var rootCmd = &cobra.Command{Use: "bookman"}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of bookman",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bookman %s", Version)
		},
	}

	rootCmd.AddCommand(bookCmd)
	rootCmd.AddCommand(collectionCmd)
	rootCmd.AddCommand(versionCmd)

	// Check for version flag in rootCmd PreRun
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("bookman %s\n", Version)
			os.Exit(0)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
