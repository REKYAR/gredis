package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func handleVersion1(conn net.Conn) {
	// Read total size (uint64)
	var totalSize uint64
	err := binary.Read(conn, binary.BigEndian, &totalSize)
	if err != nil {
		checkErr(err)
		return
	}

	// Read metadata size (uint64)
	var metadataSize uint64
	err = binary.Read(conn, binary.BigEndian, &metadataSize)
	if err != nil {
		checkErr(err)
		return
	}

	// Read metadata
	metadata := make([]byte, metadataSize)
	_, err = io.ReadFull(conn, metadata)
	if err != nil {
		checkErr(err)
		return
	}

	// Read message
	messageSize := totalSize - 16 - metadataSize // 16 bytes for two uint64 fields
	message := make([]byte, messageSize)
	_, err = io.ReadFull(conn, message)
	if err != nil {
		checkErr(err)
		return
	}

	fmt.Println("Metadata:", string(metadata))
	fmt.Println("Message:", string(message))
}
