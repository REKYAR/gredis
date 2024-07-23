package main

import (
	"fmt"
	"io"
	"net"
)

const MX_REQ_LEN = 4096
const MX_READ_LEN = 512

func main() {
	l, err := net.Listen("tcp", ":8080")
	checkErr(err)
	defer l.Close()
	fmt.Println("Listening on localhost:8080")
	for {
		conn, err := l.Accept()
		checkErr(err)
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	recvBuf := make([]byte, MX_REQ_LEN)
	tmpBuf := make([]byte, MX_READ_LEN)
	defer conn.Close()
	//for {
	for {
		n, err := conn.Read(tmpBuf)
		if err != nil {
			if err != io.EOF {
				checkErr(err)
			} else {
				break
			}
		}
		recvBuf = append(recvBuf, tmpBuf[:n]...)
	}
	print(string(recvBuf) + "\n")
	//process recvBuff here
	command_array := parse_req(string(recvBuf))
	execute_command(command_array)
	//}
}
