package main

import (
	"github.com/luism2302/gator/internal/config"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	err = cfg.SetUser("luis")
	if err != nil {
		log.Fatalf("error setting user: %v", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
}
