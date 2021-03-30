package main

import (
	"github.com/jotitan/monitor-pis/agent"
	"github.com/jotitan/monitor-pis/config"
	"log"
)

func main(){
	conf := config.NewConfig("agent-config.json")
	ag,err := agent.NewAgentRunner(conf)
	if err != nil {
		log.Fatal(err)
	}
	ag.Run()
}