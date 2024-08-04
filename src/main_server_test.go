package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
)

// MockConn is a mock implementation of net.Conn
type MockConn struct {
	net.Conn
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *MockConn) Close() error {
	return nil
}

func TestHandleConnection(t *testing.T) {
	tests := []struct {
		name            string
		protocolVersion uint64
		expectedOutput  string
	}{
		{
			name:            "Supported protocol version 1",
			protocolVersion: 1,
			expectedOutput:  "Handled version 1",
		},
		{
			name:            "Unsupported protocol version",
			protocolVersion: 2,
			expectedOutput:  "Unsupported protocol version: 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock connection
			readBuffer := new(bytes.Buffer)
			writeBuffer := new(bytes.Buffer)
			mockConn := &MockConn{
				readBuffer:  readBuffer,
				writeBuffer: writeBuffer,
			}

			// Write the protocol version to the read buffer
			err := binary.Write(readBuffer, binary.BigEndian, tt.protocolVersion)
			if err != nil {
				t.Fatalf("Failed to write protocol version to buffer: %v", err)
			}

			// Call handleConnection with the mock connection
			handleConnection(mockConn)

			// Check the output
			output := writeBuffer.String()
			if !bytes.Contains([]byte(output), []byte(tt.expectedOutput)) {
				t.Errorf("Expected output to contain %q, got %q", tt.expectedOutput, output)

			}
		})
	}
}
