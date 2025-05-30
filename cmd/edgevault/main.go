package main

import (
	"fmt"
	"os"

	"github.com/DhruvDattani1/edgevault/internal/storage"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: edgevault put <source_file>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	sourceFile := os.Args[2]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		fmt.Printf("Error: source file '%s' does not exist\n", sourceFile)
		os.Exit(1)
	}
	//Make sure the file exists ^

	masterKey := []byte("12345678901234567890123456789012") // 32 bytes

	switch cmd {
	case "put":
		err := storage.Put(sourceFile, masterKey)
		if err != nil {
			fmt.Printf("Error putting file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File stored successfully.")
	default:
		fmt.Println("Unknown command.")
	}
}
