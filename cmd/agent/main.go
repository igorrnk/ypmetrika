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
	//logger.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	config := configs.InitAgentConfig()

	agent, err := agents.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
