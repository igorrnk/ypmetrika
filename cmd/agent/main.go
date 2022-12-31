package main

import (
	"github.com/igorrnk/ypmetrika/configs"
	"github.com/igorrnk/ypmetrika/internal/agents"
	"log"
	"os"
)

func main() {
	logger := log.Default()
	logger.SetOutput(os.Stdout)

	config := configs.InitAgentConfig()

	agent, err := agents.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
