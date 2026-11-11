/*
1. Define a load balancer server which asks to listen on a port using net.Listen()
2. Start accepting a request from a client and verify if the request is reaching my server
3. Introduce 1 backend server and connect that server to the load balancer
4. Create a bookkeeping mechanism for the load balancer? (or pass through.. this is the main plumbing of the entire project)
5. Set up communication between both the connections
6. Write code for closing connections (all situations we covered)
7. Introduce a second server.. same mechanism.. but now we introduce actual load balancing..
8. Later add more features... but this is the baseline
*/

package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("listening on 127.0.0.1:8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

// this represents handling a single client connection to the load balancer
func handle(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("\n--- new connection from %s ---\n", conn.RemoteAddr())

	buf := make([]byte, 4096)

	backend_conn, err := net.Dial("tcp", "127.0.0.1:8081") // right now hardcoding the address.. with multiple servers, will add context from the connection using for assignment (what goes where)
	if err != nil {
		log.Println("error while connecting to backend:", err)
		return
	}
	defer backend_conn.Close()

	go handleBackendResponse(backend_conn, conn)
	for {
		n, err := conn.Read(buf) // read whatever bytes have arrived so far
		if err != nil {
			log.Println("client→backend: read => error while reading from connection:", err)
			return // TODO: handle connection closure and cleanup properly.. right now we simply close the connection.. introduce some sort of a graceful shutdown for the connection
		}

		fmt.Println(string(buf[:n])) // human-readable
		fmt.Printf("%q\n", buf[:n])  // same bytes, with \r\n made visible

		_, err = backend_conn.Write(buf[:n])
		if err != nil {
			log.Println("client→backend: write => error while writing to backend connection:", err)
			return
		}
	}
}

func handleBackendResponse(backend_conn net.Conn, client_conn net.Conn) {
	defer backend_conn.Close()

	buf := make([]byte, 4096)

	for {
		n, err := backend_conn.Read(buf)
		if err != nil {
			log.Println("backend→client: read => error while reading from connection:", err)
			return
		}

		fmt.Printf("%s\n", string(buf[:n]))
		fmt.Printf("%q\n", buf[:n])

		_, err = client_conn.Write(buf[:n])
		if err != nil {
			log.Println("backend→client: write => error while writing to client connection:", err)
			return
		}
	}
}
