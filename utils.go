package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
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

func printDataTable(dataHeader []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(dataHeader)
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
