package main

import (
	"fmt"
	"log"

	"crypto-watcher-backend/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Println("APPLICATION CONFIG", cfg)
}
