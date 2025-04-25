package main_test

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
)

func TestSimple(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	commands := [][]byte{
		{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}, // I 12345 101
		{0x49, 0x00, 0x00, 0x30, 0x3a, 0x00, 0x00, 0x00, 0x66}, // I 12346 102
		{0x49, 0x00, 0x00, 0x30, 0x3b, 0x00, 0x00, 0x00, 0x64}, // I 12347 100
		{0x49, 0x00, 0x00, 0xa0, 0x00, 0x00, 0x00, 0x00, 0x05}, // I 40960 5
		{0x51, 0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x40, 0x00}, // Q 12288 16384
	}

	for _, cmd := range commands {
		_, err := conn.Write(cmd)
		if err != nil {
			t.Fatalf("failed to send command: %v", err)
		}
	}

	response := make([]byte, 4)
	_, err = conn.Read(response)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	expected := make([]byte, 4)
	binary.BigEndian.PutUint32(expected, 101)

	if !bytes.Equal(response, expected) {
		t.Errorf("unexpected response: got %v, want %v", response, expected)
	}
}
