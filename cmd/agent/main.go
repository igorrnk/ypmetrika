package main

import (
	"github.com/igorrnk/ypmetrika/internal/agents"
	"log"
	"os"
)

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stdout)

	agent := agents.NewAgent()
	err := agent.FillMetrics()
	if err != nil {
		log.Fatal(err)
	}
	agent.Run()
}
