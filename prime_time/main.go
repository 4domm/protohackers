package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
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
		go handle(conn)
	}
}

type message struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}
type answer struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		var msg message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			_, _ = conn.Write([]byte("error\n"))
			return
		}
		if msg.Method != "isPrime" || msg.Number == nil {
			_, _ = conn.Write([]byte("error\n"))
			return
		}
		var ans answer
		if math.Trunc(*msg.Number) == *msg.Number {
			ans = answer{Method: "isPrime", Prime: isPrime(int64(*msg.Number))}
		} else {
			ans = answer{Method: "isPrime", Prime: false}
		}
		buff, _ := json.Marshal(ans)
		_, _ = conn.Write(append(buff, '\n'))
	}
}

func isPrime(n int64) bool {
	if n <= 1 {
		return false
	}
	for i := int64(2); i <= int64(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
