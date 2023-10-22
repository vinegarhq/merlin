/*
TODO:
Get CSV writing complete
Detect Duplicates
Parse the survey fields
*/

package main

import (
	"fmt"
	"github.com/vinegarhq/merlin/internal"
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

	config, err := internal.LoadConfiguration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	err = internal.BeginListener(config)
	if err != nil {
		log.Fatal(err)
	}
}
