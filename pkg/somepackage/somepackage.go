package somepackage

import (
	"fmt"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func ExampleFunction() {
	err := someOperation()
	CheckError(err)
}

// someOperation is a dummy function to simulate an operation that may return an error
func someOperation() error {
	// Simulate an error
	return fmt.Errorf("an example error occurred")
}
