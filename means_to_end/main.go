package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
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
	type entry struct {
		timestamp int32
		price     int32
	}
	prices := make([]entry, 0)
	for {
		message := make([]byte, 9)
		_, err := io.ReadFull(conn, message)
		if err != nil {
			return
		}
		if message[0] == byte('I') {
			timestamp := int32(binary.BigEndian.Uint32(message[1:5]))
			price := int32(binary.BigEndian.Uint32(message[5:]))
			prices = append(prices, entry{timestamp: timestamp, price: price})
			fmt.Println("insert done")
		} else if message[0] == byte('Q') {
			sort.Slice(prices, func(i, j int) bool {
				return prices[i].timestamp < prices[j].timestamp
			})
			minTime := int32(binary.BigEndian.Uint32(message[1:5]))
			maxTime := int32(binary.BigEndian.Uint32(message[5:]))
			sum := int64(0)
			count := int64(0)
			ans := int32(0)
			for _, val := range prices {
				if minTime <= val.timestamp && val.timestamp <= maxTime {
					sum += int64(val.price)
					count++
				}
				if count == 0 {
					ans = 0
				} else {
					ans = int32(sum / count)
				}
			}
			if maxTime < minTime {
				ans = 0
			}
			_ = binary.Write(conn, binary.BigEndian, ans)
		}
	}

}
