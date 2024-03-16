package helper

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestReadFromFile(t *testing.T) {
	// Create a temporary file for testing
	tmpfile := "testfile.txt"
	content := "line 1\nline 2\nline 3\n"
	err := os.WriteFile(tmpfile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer os.Remove(tmpfile)

	expected := []string{"line 1\n", "line 2\n", "line 3\n"}

	// Call the function
	result, err := ReadFromFile(tmpfile)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}

	// Check if the result matches the expected output
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result doesn't match expected. Got: %v, Expected: %v", result, expected)
	}
}

func TestReadFromFileError(t *testing.T) {
	// Call the function with a non-existent file
	_, err := ReadFromFile("nonexistent.txt")

	// Check if an error is returned
	if err == nil {
		t.Error("expected error but got none")
	} else {
		fmt.Println("Error:", err)
	}
}
