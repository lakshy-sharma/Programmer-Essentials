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
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"text/tabwriter"
	"time"

	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"
	"github.com/tatsushid/go-fastping"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "A tool to discover online hosts in your network.",
	Long: `The discover tool allows you to scan all hosts inside a network and check if they are online or not.
	It is capable of mapping IPs to their hostnames, making it easier to find a rogue raspberry pi ;)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		scanResults := discoverHosts()
		writer := tabwriter.NewWriter(os.Stdout, 1, 8, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Fprintln(writer, "\nScan Complete")
		fmt.Fprintln(writer, "--------------------------------------------")
		fmt.Fprintln(writer, "IP Address\tState\tHostname\tResponse Time")
		fmt.Fprintln(writer, "--------------------------------------------")
		for _, result := range scanResults {
			if result.state == "Up" {
				fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", result.ipaddress, result.state, result.hostname, result.responseTime)
			}
		}
		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
	discoverCmd.Flags().StringVarP(&networkCidr, "cidr", "c", "192.168.0.0/24", "The CIDR notation of the network you want to scan.")
}

var (
	networkCidr string
)

type ipData struct {
	ipaddress    string
	state        string
	hostname     []string
	responseTime time.Duration
}

type pingResult struct {
	ipAddress    *net.IPAddr
	ipState      string
	responseTime time.Duration
}

const PING_TIMER = 10

// This function captures the IP addresses which must be scanned based on shown network CIDR.
// It returns a slice of valid net.IP instances.
func convertToIPs() []net.IP {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(networkCidr)
	if err != nil {
		log.Fatal(err)
	}

	// convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32 and store them in a slice.
	ipStore := []net.IP{}
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ipStore = append(ipStore, ip)
	}
	return ipStore
}

// This function setups the ping control mechanism
func pingController(ipsToScan []net.IP, pingResultChannel chan pingResult, finishChannel chan string) {
	// Setup the IP pinger.
	p := fastping.NewPinger()
	p.MaxRTT = time.Second*PING_TIMER + 1
	var pingOutput pingResult

	for _, ip := range ipsToScan {
		// Resolve the ip addresses and add them to a IP address pinger.
		resolvedAddress, err := net.ResolveIPAddr("ip4:icmp", ip.String())
		if err != nil {
			os.Exit(1)
		}
		p.AddIPAddr(resolvedAddress)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		pingOutput.ipAddress = addr
		pingOutput.ipState = "Up"
		pingOutput.responseTime = rtt
		pingResultChannel <- pingOutput
	}

	// Once the final time passes stop the program.
	p.OnIdle = func() {
		finishChannel <- "Completed"
	}
	err := p.Run()
	if err != nil {
		panic(err)
	}
}

// This function is responsible for sending the ping packets to the IP hosts and collecting other details for each IP.
func formatResults(pingResultChannel chan pingResult, ipDataResults chan ipData, ipsToScan []net.IP, finishChannel chan string) {

	// Wait for each ping result and perform a IP lookup.
	for {
		select {
		case finishMessage := <-finishChannel:
			if finishMessage == "Completed" {
				finishChannel <- "Finish"
			}

		case pingOutput := <-pingResultChannel:
			var ipDetails ipData

			// Capture the IP data.

			ipDetails.ipaddress = pingOutput.ipAddress.String()
			ipDetails.state = pingOutput.ipState
			ipDetails.responseTime = pingOutput.responseTime

			// Perform a Name Lookup.
			lookup, err := net.LookupAddr(pingOutput.ipAddress.String())
			if err == nil {
				ipDetails.hostname = lookup
			} else {
				ipDetails.hostname = []string{"N/A"}
			}
			ipDataResults <- ipDetails
		}
	}
}

// This function is the main driver function responsible for controlling the flow of discovery.
func discoverHosts() []ipData {
	ipsToScan := convertToIPs()
	pingResultChannel := make(chan pingResult)
	ipDataResults := make(chan ipData)
	finishChannel := make(chan string)

	var scanResults []ipData
	go pingController(ipsToScan, pingResultChannel, finishChannel)
	go formatResults(pingResultChannel, ipDataResults, ipsToScan, finishChannel)

	fmt.Printf("This scan will run for %d Seconds to find LAN peers.\n", PING_TIMER)
	// Show a simple Progress Bar.
	go func() {
		bar := progressbar.New(PING_TIMER)
		for i := 1; i <= PING_TIMER; i++ {
			bar.Add(1)
			time.Sleep(1 * time.Second)
		}
		fmt.Print("\n")
	}()

	for {
		select {
		case ipResults := <-ipDataResults:
			scanResults = append(scanResults, ipResults)
		case finishMessage := <-finishChannel:
			if finishMessage == "Finish" {
				return scanResults
			}
		}
	}
}
