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
	"log"
	"net"
	"strconv"
	"sync"
)

func processTCPClient(clientConnection net.Conn, replyMessage string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Received a new client connection.")
	for {
		message, err := bufio.NewReader(clientConnection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		// Exit when client sends EOF
		if message == "EOF" {
			return
		}

		// Display the received message.
		fmt.Println("<- ", string(message))

		// Client asked for a echo server then reformat the message and send back.
		// Else simply send back the required response.
		if replyMessage == "ECHO" {
			clientConnection.Write([]byte(fmt.Sprintf("Echo: %s", string(message))))
		} else {
			clientConnection.Write([]byte(string(replyMessage + "\n")))
		}
	}
}

// This function starts a TCP server with provided port and reply mechanism.
func ServeTCP(portNumber int, replyMessage string) {
	wg := new(sync.WaitGroup)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(portNumber))
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	log.Printf("TCP server started.\nPort: %d\nReply: %s\n", portNumber, replyMessage)

	for {
		clientConnection, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go processTCPClient(clientConnection, replyMessage, wg)
	}
}
