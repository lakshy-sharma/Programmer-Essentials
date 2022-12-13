/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "matrix",
	Short: "Matrix is a network oriented tool which provides you several network related functionalities.",
	Long: `The Matrix application provides you several network related functionalities for debugging and troubleshooting.

	A few have been listed below 
	1. Scan a target host for any open ports.
	2. Scan a network for available hosts.
	3. Create a TCP server for testing TCP based applications.
	4. Create a gRPC based server for testing gRPC based applications.
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
