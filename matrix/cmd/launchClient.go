/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"matrix/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	serverPort int
	serverHost string
)

// launchClientCmd represents the launchTestClient command
var launchClientCmd = &cobra.Command{
	Use:   "launchClient",
	Short: "Launch a interactive client to test server responses.",
	Long: `This command launches a client for a websocket, TCP or a gRPC server.
	It opens a interactive prompt and allows users to send customized messages to the server and test its output.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.TcpClient(serverPort, serverHost)
	},
}

func init() {
	rootCmd.AddCommand(launchClientCmd)
	launchClientCmd.Flags().IntVarP(&serverPort, "serverport", "p", 5000, "The port number on which your server is active.")
	launchClientCmd.Flags().StringVarP(&serverHost, "server", "s", "localhost", "The address where your server is active.")
}