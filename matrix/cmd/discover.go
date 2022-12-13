/*
Copyright © 2022 Lakshy Sharma lakshy1106@protonmail.com
*/
package cmd

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "A tool to discover online hosts in your network.",
	Long: `The discover tool allows you to scan all hosts inside a network and check if they are online or not.
	It is capable of mapping IPs to their hostnames, making it easier to find a rogue raspberry pi ;)
	`,
	Run: func(cmd *cobra.Command, args []string) {
		discoverHosts()
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
	ipaddress string
	State     string
	hostname  []string
}

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

func captureIpDetails(ip net.IP) ipData {
	var ipDetails ipData
	// Set the Ip address
	ipDetails.ipaddress = ip.String()

	// Lookup the name of a IP
	lookup, err := net.LookupAddr(ip.String())
	if err == nil {
		ipDetails.hostname = lookup
	}
	return ipDetails
}

func discoverHosts() {
	ipsToScan := convertToIPs()

	for _, ip := range ipsToScan {
		ipDetails := captureIpDetails(ip)
		fmt.Printf("%s\t%s\t%s\n", ipDetails.ipaddress, ipDetails.State, ipDetails.hostname)
	}
}
