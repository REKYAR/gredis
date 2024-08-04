package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	MX_REQ_LEN  = 4096
	MX_READ_LEN = 512
	host        = "localhost"
	port        = 8080
)

var supported_protocol_versions = [...]uint64{1}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	checkErr(err)
	defer l.Close()
	fmt.Println("Listening on", fmt.Sprintf("%s:%d", host, port))
	for {
		conn, err := l.Accept()
		checkErr(err)
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	// Read protocol version (uint64)
	var version uint64
	err := binary.Read(conn, binary.BigEndian, &version)
	if err != nil {
		checkErr(err)
		return
	}
	handled := false
	for _, v := range supported_protocol_versions {
		if version == v {
			handled = true
			fmt.Println("Handled version", version)
			conn.Write([]byte(fmt.Sprintf("Handled version %d", version)))
			break
		}
	}
	if !handled {
		fmt.Println("Unsupported protocol version: ", version)
		conn.Write([]byte(fmt.Sprintf("Unsupported protocol version: %d", version)))
	}

	//TODO RESPOND TO CLIENT
}
