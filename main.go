package main

import (
	"cobra/cmd"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(cmd.WordCmd, cmd.DelKeyCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute error: %v", err)
	}
}
