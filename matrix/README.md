# Matrix

The matrix is a collection of networking tools written in Go using cobra CLi framework.\n
The focus is on increasing the performance of available utilities and create more to support developers.\n
It currently allows you to perform the following actions.

## Features
1. Scan a particular host for open ports.
2. Scan a network for hosts that are active. (This feature needs superuser access)

## Example
1. Find open ports on a host: <i>matrix portScan [IP address to scan]</i>
2. Find active hosts on a network: <i>matrix hostScan [Network CIDR to scan]</i>

## TODO
1. Add feature of hosting a simple TCP server and a websocket server for testing.
2. Add feature for creating network packets for testing high speed networks.
