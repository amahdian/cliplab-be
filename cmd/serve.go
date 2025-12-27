package cmd

import (
	"log"

	"github.com/amahdian/cliplab-be/global/env"
	"github.com/amahdian/cliplab-be/pkg/logger"
	server2 "github.com/amahdian/cliplab-be/server"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a REST server",
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	envs, err := env.Load("")
	if err != nil {
		log.Fatalf("failed to load env variables: %v", err)
	}

	server, err := server2.NewServer(envs)
	if err != nil {
		logger.Fatalf("failed to load server, err: %v", err)
	}

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
