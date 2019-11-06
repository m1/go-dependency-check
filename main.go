package main

import (
	"log"

	"github.com/m1/go-dependency-check/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
