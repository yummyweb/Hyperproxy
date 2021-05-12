package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
    Use:   "hyperproxy",
    Short: "A reverse proxy for the experts.",
    Long: `Hyperproxy is a production-ready reverse proxy which allows 
	you to setup and proxy urls through the CLI.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}