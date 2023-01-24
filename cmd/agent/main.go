package main

import (
	"github.com/igorrnk/ypmetrika/internal/agents"
	"github.com/igorrnk/ypmetrika/internal/configs"
	"log"
	"os"
)

func main() {
	//logFile, _ := os.OpenFile("./log/agentLog.log", os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0644)
	log.SetOutput(os.Stdout)
	//log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	config := configs.InitAgentConfig()

	agent, err := agents.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
