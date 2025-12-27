package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cliplab-be",
	Short: "cliplab",
	Long:  "cliplab Backend Application",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
