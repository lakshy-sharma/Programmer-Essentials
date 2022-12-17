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
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

/*
The scan command utilities.
*/

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perform a TCP connect scan on a given network host.",
	Long: `A TCP connect scan simply connects with a port on a remote server and checks if it open or not.
	This scan has been implemented in parallel fashion to make it quick.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting a TCP connect scan")

		// Gather TCP scan results and clean them.
		tcpResults := scanHost(hostname, startPort, endPort)
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

/*
Code variables.
*/

// Cant let you run wild with this thing. Can I?
// Important. Don't be a douche and turn this setting up wildly,
// you might accidentally launch a mild DOS attack.
// Trust me Bigger is not always better for performance.
const SANITY_LIMIT = 50

var (
	hostname  string
	startPort int
	endPort   int
)

type ScanResult struct {
	Port    int
	State   string
	Service string
}

/*
Helping Functions
*/
// This function scans a port on a particular host and returns the result in a struct.
func scanPort(protocol string, hostname string, port int, portResultChannel chan ScanResult) {
	result := ScanResult{Port: port, Service: protocol}
	address := hostname + ":" + strconv.Itoa(port)
	connect, err := net.DialTimeout(protocol, address, 10*time.Second)
	if err != nil {
		result.State = "Closed"
		portResultChannel <- result
		return
	}
	defer connect.Close()
	result.State = "Open"
	portResultChannel <- result
}

func resultCollector(startPort int, endPort int, portResultChannel chan ScanResult, resultCaptureChannel chan []ScanResult) {
	var results []ScanResult

	for port := startPort; port <= endPort; port++ {
		scanOutput := <-portResultChannel
		results = append(results, scanOutput)
	}

	// Once all outputs have been collected send them back to our main thread.
	resultCaptureChannel <- results
	close(resultCaptureChannel)
	close(portResultChannel)
}

/*
Main scan controller function.
This function spawns multiple goroutines to scan the ports on a host and then waits for them to finish before moving ahead.
*/
func scanHost(hostname string, startPort int, endPort int) []ScanResult {
	speedlimitChannel := make(chan struct{}, SANITY_LIMIT)
	portResultChannel := make(chan ScanResult)
	resultCaptureChannel := make(chan []ScanResult)
	wg := sync.WaitGroup{}
	defer wg.Wait()

	// Start a receiver for capturing the outputs of our scan.
	wg.Add(1)
	go func() {
		defer wg.Done()
		resultCollector(startPort, endPort, portResultChannel, resultCaptureChannel)
	}()

	// Scan Ports asynchronously.
	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		speedlimitChannel <- struct{}{}
		go func(hostname string, port int, returnChannel chan ScanResult) {
			defer wg.Done()
			scanPort("tcp", hostname, port, returnChannel)
			<-speedlimitChannel
		}(hostname, port, portResultChannel)
	}

	// Capture and clean the TCP scan results.
	finalResult := <-resultCaptureChannel
	sort.SliceStable(finalResult, func(i, j int) bool {
		return finalResult[i].Port < finalResult[j].Port
	})
	return finalResult
}
