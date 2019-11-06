package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/m1/go-dependency-check/api"
)

var (
	serveCmd = cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}
)

func serve() {
	cfg := api.Config{
		Redis: os.Getenv("REDIS"),
		Port:  os.Getenv("PORT"),
	}

	a := api.New(cfg)
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
