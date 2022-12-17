package utils

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/schollz/progressbar"
	"github.com/tatsushid/go-fastping"
)

type ipData struct {
	Ipaddress    string
	State        string
	Hostname     []string
	ResponseTime time.Duration
}

type pingResult struct {
	ipAddress    *net.IPAddr
	ipState      string
	responseTime time.Duration
}

// This function takes a network CIDR and processes it.
// The output is a list of IPs contained within the network CIDR.
func convertToIPs(networkCidr string) []net.IP {
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

// This function is responsible for sending ping packets to all the machines within the network.
// Once we start receiving replies we send them to the output formatter for further processing.
func pingSender(ipsToScan []net.IP, pingResultChannel chan pingResult, finishChannel chan string, pingTimer int) {
	// Setup the IP pinger.
	p := fastping.NewPinger()
	p.MaxRTT = time.Second*time.Duration(pingTimer) + 1
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

// This function receives the results from the ping Sender and formats them with more information to make it presentable.
func outputFormatter(pingResultChannel chan pingResult, ipDataResults chan ipData, ipsToScan []net.IP, finishChannel chan string) {

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

			ipDetails.Ipaddress = pingOutput.ipAddress.String()
			ipDetails.State = pingOutput.ipState
			ipDetails.ResponseTime = pingOutput.responseTime

			// Perform a Name Lookup.
			lookup, err := net.LookupAddr(pingOutput.ipAddress.String())
			if err == nil {
				ipDetails.Hostname = lookup
			} else {
				ipDetails.Hostname = []string{"N/A"}
			}
			ipDataResults <- ipDetails
		}
	}
}

// A function to show a simple progress bar.
func progressBar(pingTimer int) {
	bar := progressbar.New(pingTimer)
	for i := 1; i <= pingTimer; i++ {
		bar.Add(1)
		time.Sleep(1 * time.Second)
	}
	fmt.Print("\n")
}

// This function is the control function which controls how the hosts are discovered within the network.
func DiscoverHosts(networkCidr string, pingTimer int) []ipData {
	// Setting up the variables and the channels for communication between the threads.
	ipsToScan := convertToIPs(networkCidr)
	pingResultChannel := make(chan pingResult)
	ipDataResults := make(chan ipData)
	finishChannel := make(chan string)
	var scanResults []ipData

	// Display welcome message and progress bar.
	fmt.Printf("This scan will run for %d Seconds to find LAN peers.\n", pingTimer)
	go progressBar(pingTimer)

	// Start a ping sender and reply receiver.
	go pingSender(ipsToScan, pingResultChannel, finishChannel, pingTimer)
	go outputFormatter(pingResultChannel, ipDataResults, ipsToScan, finishChannel)

	// Capture all the replies as the outputFormatter sends them.
	// Close when the formatter says it is done.
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
