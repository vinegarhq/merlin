package main

import (
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: merlin [path to configuration JSON]")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	config, err := LoadConfiguration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Loaded configuration successfully.")

	log.Println("Starting HTTP listener...")
	err = BeginListener(config)
	if err != nil {
		log.Fatal(err)
	}
}
