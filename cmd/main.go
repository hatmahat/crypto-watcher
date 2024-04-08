package main

import (
	"log"

	"crypto-watcher-backend/cmd/worker"
	"crypto-watcher-backend/internal/config"

	"github.com/spf13/cobra"
)

func main() {
	cfg, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	cmd := &cobra.Command{
		Use:   "crypto-watcher",
		Short: "Watch & Alert Crypto Prices",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:              "worker",
			Short:            "Worker Server",
			Long:             "a worker server that will run all scheduler registered to the worker",
			TraverseChildren: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				worker.Start(cfg)
				return nil
			},
		},
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
