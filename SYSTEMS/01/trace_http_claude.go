// Stage 0, experiment 3: a tiny TCP server that shows you the raw bytes
// a browser actually sends. No HTTP library, just a socket.
//
// Run:   go run main.go
// Then open http://localhost:8080 in your browser.
package main

import (
        "fmt"
        "log"
        "net"
)

func main() {
        // Ask the OS for a listening TCP socket on port 8080.
        ln, err := net.Listen("tcp", ":8080")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("listening on :8080")

        for {
                // Block until a client completes the TCP handshake with us.
                conn, err := ln.Accept()
                if err != nil {
                        log.Println("accept:", err)
                        continue
                }
                go handle(conn) // one goroutine per connection
        }
}

func handle(conn net.Conn) {
        defer conn.Close()
        fmt.Printf("\n--- new connection from %s ---\n", conn.RemoteAddr())

        buf := make([]byte, 4096)
        n, err := conn.Read(buf) // read whatever bytes have arrived so far
        if err != nil {
                log.Println("read:", err)
                return
        }

        fmt.Println(string(buf[:n])) // human-readable
        fmt.Printf("%q\n", buf[:n])  // same bytes, with \r\n made visible

        // STEP B: uncomment the line below, re-run, and refresh the browser.
        // You are writing an HTTP response by hand.
        // conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 6\r\n\r\nhello\n"))
}