package main

import (
	"bufio"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (env *runEnv) GetUserTextInput() string {
	reader := bufio.NewReader(env.userInput)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// ReadString keeps delim, remove it here
	text = strings.TrimSuffix(text, "\n")
	return text
}
