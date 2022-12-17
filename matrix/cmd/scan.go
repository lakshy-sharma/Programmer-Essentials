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
	"fmt"
	"matrix/pkg/utils"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	hostname  string
	startPort int
	endPort   int
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a TCP connect scan on a given network host.",
	Long: `A TCP connect scan simply connects with a port on a remote server and checks if it open or not.
	This scan has been implemented in parallel fashion to make it quick.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Gather TCP scan results and clean them.
		tcpResults := utils.ScanHostPorts(hostname, startPort, endPort)
		sort.SliceStable(tcpResults, func(i, j int) bool {
			return tcpResults[i].Port < tcpResults[j].Port
		})

		// Print the results in a clean fashion.
		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(writer, "Port\tState\tService")
		fmt.Fprintln(writer, "----------------------")
		for _, result := range tcpResults {
			if result.State == "Open" {
				fmt.Fprintf(writer, "%d\t%s\t%s\n", result.Port, result.State, result.Service)
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&hostname, "hostname", "H", "localhost", "The host you want to scan.")
	scanCmd.Flags().IntVarP(&startPort, "start_port", "s", 1, "Start number of the port you want to scan.")
	scanCmd.Flags().IntVarP(&endPort, "end_port", "e", 1024, "The port number you want to stop scanning at.")
}
