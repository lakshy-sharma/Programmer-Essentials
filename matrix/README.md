# Matrix

The matrix is a collection of networking tools written in Go using cobra CLi framework.\n
The focus is on increasing the performance of available utilities and create more to support developers.\n
It currently allows you to perform the following actions.

## Features
1. Scan a particular host for open ports.
2. Scan a network for hosts that are active. (This feature needs superuser access)
3. Launch a test TCP server for testing TCP clients.
4. Launch a test TCP client for testing TCP servers.

## Example
1. Find open ports on a host: <i>matrix portScan -H [IP address to scan] -s [Start port] -e [End port]</i>
2. Find active hosts on a network: <i>matrix hostScan -c [Network CIDR to scan] -t [Time for Ping reply]</i>

## TODO
1. Extend the TCP server and client to include Websockets and gRPC.
2. Add feature for creating network packets for testing high speed networks.
