package main

import (
	"log"

	"github.com/zalhui/calc_golang/config"
	"github.com/zalhui/calc_golang/internal/agent/worker"
)

func main() {
	cfg := config.LoadConfig()

	for i := 0; i < cfg.ComputingPower; i++ {
		go worker.StartWorker()
	}

	log.Printf("Agent started with %d workers\n", cfg.ComputingPower)
	select {}

	//defer log.Println("Shutting down agent...")
	//defer log.Println("Agent stopped")
	//cancel()
}
