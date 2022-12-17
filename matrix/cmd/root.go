/*
Copyright © [2022] [Lakshy Sharma] <lakshy.sharma@protonmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	Long: `
	WELCOME TO THE MATRIX.
	The Matrix application provides developers with several network related functionalities.
	It aims to serve as the swiss army knife of network application testing and debugging.

	Features: Currently available.
	1. Scan a target host for any open ports.
	2. Scan a network for available hosts.

	Upcoming Features: 
	1. A Simple TCP/WEbsocket or gRPC server for testing your peer to peer clients.
	2. A high speed packet generator for testing networks.
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
