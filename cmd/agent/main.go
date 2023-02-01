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

	config, err := configs.InitAgentConfig()
	if err != nil {
		log.Fatal(err)
	}
	agent, err2 := agents.NewAgent(config)
	if err2 != nil {
		log.Fatal(err2)
	}
	if err := agent.Run(); err != nil {
		log.Fatal(err)
	}
}
