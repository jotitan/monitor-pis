package main

import (
	"github.com/jotitan/monitor-pis/config"
	"github.com/jotitan/monitor-pis/monitor_server"
	"log"
)

func main() {
	c, err := config.NewMonitoringConfig("monitoring.json")
	log.Println("Use monitoring.json configuration")
	if err != nil {
		log.Fatal(err)
	}
	ms := monitor_server.NewMonitoringServer(*c)
	ms.Start()
}
