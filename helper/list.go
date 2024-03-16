package helper

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

func RemoveColorFromList(list1 []string, color string) []string {
	var out []string
	for _, item := range list1 {
		if strings.Contains(item, color) {
			continue
		}
		out = append(out, item)
	}
	return out
}

func MergeLists(list1, list2 []string) []string {
	// Convert list1 to a map for efficient lookup
	existing := make(map[string]bool)
	for _, item := range list1 {
		// Extract X and Y values
		parts := strings.Fields(item)
		if len(parts) >= 3 { // Ensure there are at least three parts
			key := parts[1] + " " + parts[2] // Concatenate X and Y values
			existing[key] = true
		}
	}

	// Append all items from list1
	finalList := make([]string, len(list1))
	copy(finalList, list1)

	// Append items from list2 that are not in list1
	for _, item := range list2 {
		// Extract X and Y values
		parts := strings.Fields(item)
		if len(parts) >= 3 { // Ensure there are at least three parts
			key := parts[1] + " " + parts[2] // Concatenate X and Y values
			if !existing[key] {
				finalList = append(finalList, item)
			}
		}
	}

	return finalList
}

func PacketizeList(input []string, maxSize int) []string {
	var result []string
	var currentString string

	for _, str := range input {
		// If adding the next string exceeds maxSize, start a new list
		if len(currentString)+len(str) > maxSize {
			result = append(result, currentString)
			currentString = ""
		}
		currentString += str
	}

	// Append the last remaining list
	if len(currentString) > 0 {
		result = append(result, currentString)
	}

	return result
}

func ShuffleStrings(slice []string) []string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(slice)
	shuffled := make([]string, n)
	copy(shuffled, slice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

// SplitStringList splits a list of strings into 'numParts' parts.
func SplitList(input []string, numParts int) [][]string {
	// Initialize the resulting parts
	parts := make([][]string, numParts)
	// Calculate the number of elements per part
	elementsPerPart := len(input) / numParts
	remainder := len(input) % numParts
	// Start index for iterating over input list
	start := 0
	// Iterate over each part
	for i := range parts {
		// Calculate the end index for this part
		end := start + elementsPerPart
		// Distribute the remainder elements equally among the first parts
		if i < remainder {
			end++
		}
		// Assign elements to this part
		parts[i] = input[start:end]
		// Move the start index for the next part
		start = end
	}
	return parts
}

func GenerateRandomColor() string {
	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)
	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}

func PrintByteSize(input []string) {

	// Calculate the total byte size
	totalByteSize := 0
	for _, str := range input {
		totalByteSize += len(str)
	}

	// If you want to include the overhead of the slice header
	// which is 24 bytes on a 64-bit system
	totalByteSize += len(input) * int(unsafe.Sizeof(input[0]))

	// Convert bytes to megabytes
	totalGB := float64(totalByteSize) / (1024 * 1024)

	fmt.Printf("Total size of input: %.2f MB\n", totalGB)
}
