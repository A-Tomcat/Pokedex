package main

import (
	"strings"
)

func cleanInput(text string) []string {
	output := strings.Fields(strings.ToLower(text))
	return output
}
