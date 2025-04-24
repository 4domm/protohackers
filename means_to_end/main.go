package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("error accepting %e", err)
			continue
		}
		go echo(conn)
	}
}

func echo(conn net.Conn) {
	defer conn.Close()

}
