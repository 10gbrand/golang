package main

import (
	"log"
)

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Kunde inte läsa konfigurationsfil: %v", err)
	}

	if err := RunMerge(config); err != nil {
		log.Fatalf("Fel vid körning: %v", err)
	}

	if err := SwapSources(config, "styles/def/sourses.json"); err != nil {
		log.Fatalf("Fel vid byte av sources: %v", err)
	}
}
