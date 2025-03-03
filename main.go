package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	return strings.Fields(lowerText)
}

func main() {
	fmt.Printf("Hello, World!")
}
