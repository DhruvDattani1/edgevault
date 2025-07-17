package main

import (
	"fmt"
	"os"

	"github.com/DhruvDattani1/edgevault/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  edgevault put <source_file>")
		fmt.Println("  edgevault get <object_name> <dest_file>")
		fmt.Println("  edgevault delete <object_name>")
		fmt.Println("  edgevault list")
		os.Exit(1)
	}

	cmd := os.Args[1]
	masterKey := []byte("12345678901234567890123456789012") // 32 bytes

	switch cmd {
	case "put":
		if len(os.Args) < 3 {
			fmt.Println("Usage: edgevault put <source_file>")
			os.Exit(1)
		}
		sourceFile := os.Args[2]
		if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
			fmt.Printf("Error: source file '%s' does not exist\n", sourceFile)
			os.Exit(1)
		}
		err := storage.Put(sourceFile, masterKey)
		if err != nil {
			fmt.Printf("Error putting file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File stored successfully.")

	case "get":
		if len(os.Args) < 4 {
			fmt.Println("Usage: edgevault get <object_name> <dest_file>")
			os.Exit(1)
		}
		objectName := os.Args[2]
		destPath := os.Args[3]
		err := storage.Get(objectName, destPath, masterKey)
		if err != nil {
			fmt.Printf("Error getting file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("File decrypted successfully.")

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: edgevault delete <object_name>")
			os.Exit(1)
		}
		objectName := os.Args[2]
		err := storage.Delete(objectName)
		if err != nil {
			fmt.Printf("Error deleting object: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Object deleted successfully.")

	case "list":
		objects, err := storage.List()
		if err != nil {
			fmt.Printf("Error listing objects: %v\n", err)
			os.Exit(1)
		}
		if len(objects) == 0 {
			fmt.Println("Vault is empty.")
			return
		}
		fmt.Println("Stored objects:")
		for _, obj := range objects {
			fmt.Printf("  - %s (%d bytes, modified %s)\n", obj.Name, obj.Size, obj.Modified.Format("Jan 02 15:04"))
		}

	default:
		fmt.Println("Unknown command.")
		os.Exit(1)
	}
}
