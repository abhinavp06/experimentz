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
	"errors"
	"io"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listening on 127.0.0.1:8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handle(conn)
	}
}

// this represents handling a single client connection to the load balancer
func handle(conn net.Conn) {
	id := conn.RemoteAddr().String()
	defer log.Printf("[%s] handle: exited, its defers ran", id)
	defer conn.Close()

	log.Printf("[%s] new connection", id)

	buf := make([]byte, 4096)

	backend_conn, err := net.Dial("tcp", "127.0.0.1:8081") // right now hardcoding the address.. with multiple servers, will add context from the connection using for assignment (what goes where)
	if err != nil {
		log.Printf("[%s] dial backend: %v", id, err)
		return
	}
	defer backend_conn.Close()
	log.Printf("[%s] dialed backend: %s -> %s", id, backend_conn.LocalAddr(), backend_conn.RemoteAddr())

	go handleBackendResponse(backend_conn, conn, id)

	for {
		n, err := conn.Read(buf) // read whatever bytes have arrived so far
		if err != nil {
			logReadErr(id, "client→backend", err)
			return // TODO: handle connection closure and cleanup properly.. right now we simply close the connection.. introduce some sort of a graceful shutdown for the connection
		}
		log.Printf("[%s] client→backend: %d bytes", id, n)

		if _, err := backend_conn.Write(buf[:n]); err != nil {
			log.Printf("[%s] client→backend: write to backend: %v", id, err)
			return
		}
	}
}

func handleBackendResponse(backend_conn net.Conn, client_conn net.Conn, id string) {
	defer log.Printf("[%s] backend→client: exited, its defers ran", id)
	defer backend_conn.Close()
	defer client_conn.Close()

	buf := make([]byte, 4096)

	for {
		n, err := backend_conn.Read(buf)
		if err != nil {
			logReadErr(id, "backend→client", err)
			return
		}
		log.Printf("[%s] backend→client: %d bytes", id, n)

		if _, err := client_conn.Write(buf[:n]); err != nil {
			log.Printf("[%s] backend→client: write to client: %v", id, err)
			return
		}
	}
}

func logReadErr(id, dir string, err error) {
	switch {
	case errors.Is(err, io.EOF):
		log.Printf("[%s] %s: read EOF (peer sent FIN)", id, dir)
	case errors.Is(err, net.ErrClosed):
		log.Printf("[%s] %s: read on a conn our own code closed", id, dir)
	default:
		log.Printf("[%s] %s: read error: %v", id, dir, err)
	}
}