package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8081", "address to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("backend listening on", *addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handle(conn, *addr)
	}
}

func handle(conn net.Conn, name string) {
	defer conn.Close()
	log.Println("new connection from", conn.RemoteAddr())

	reader := bufio.NewReader(conn) // buffers across partial reads
	for {
		line, err := reader.ReadString('\n') // blocks until a full line (incl. \n) or an error
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("client closed:", conn.RemoteAddr()) // normal goodbye
			} else {
				log.Println("read:", err) // something actually broke
			}
			return
		}

		log.Printf("got %q", line)

		if _, err := fmt.Fprintf(conn, "[%s] %s", name, line); err != nil {
			log.Println("write:", err)
			return
		}
	}
}
