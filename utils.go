package main

import (
	"bufio"
	"os"
)

func GetUserTextInput() []string {
	userInput := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		text := scanner.Text()
		if len(text) != 0 {
			userInput = append(userInput, text)
		} else {
			break
		}

	}
	return userInput
}
