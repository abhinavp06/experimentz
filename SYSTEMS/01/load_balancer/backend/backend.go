package main

import (
        "fmt"
        "log"
        "net"
)

func main() {
        ln, err := net.Listen("tcp", "127.0.0.1:8081")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println("listening on 127.0.0.1:8081")

        for {
                conn, err := ln.Accept()
                if err != nil {
                        log.Println("accept:", err)
                        continue
                }
                go handle(conn)
        }
}

func handle(conn net.Conn) {
        defer conn.Close()
        fmt.Printf("\n--- new connection from %s ---\n", conn.RemoteAddr())

        buf := make([]byte, 4096)
        _, err := conn.Read(buf) // read whatever bytes have arrived so far
        if err != nil {
                log.Println("read:", err)
                return
        }

        conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 6\r\n\r\nhello\n"))
}