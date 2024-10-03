package main

import (
	"log"

	"github.com/Layer-Edge/bitcoin-da/config"
	"github.com/Layer-Edge/bitcoin-da/da" // Adjust import path as needed.
)

func main() {
	cfg := loadConfig() // Assume loadConfig is a function that reads your config file.

	go da.RawBlockSubscriber(cfg) // Start the raw block subscriber.
	go da.HashBlockSubscriber(cfg) // Start the hash block subscriber.

	select {} // Keep the main function running.
}

func loadConfig() *config.Config {
	var cfg config.Config
	
	// Load your configuration here (e.g., from a JSON or YAML file).
	return &cfg
}
