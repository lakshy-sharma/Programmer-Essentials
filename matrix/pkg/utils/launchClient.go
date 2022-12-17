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
package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

// This function starts a TCP client which connects with the server and allows users to connect test server responses.
func TcpClient(serverPort int, serverHost string) {
	serverConnection, err := net.Dial("tcp", serverHost+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Type in the message you want to send to the server.")
	for {
		// Reading the input from the user to send back to the server.
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintln(serverConnection, text)

		// Capturing the output from the server and displaying it.
		message, _ := bufio.NewReader(serverConnection).ReadString('\n')
		fmt.Print("-> " + message)
	}
}
