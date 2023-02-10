package functional_options

import (
	"fmt"
	"testing"
)

func TestMainer(t *testing.T) {
	server := NewServer(
		WithPort(8080),
		WithProtocol("http"),
	)

	fmt.Printf("Server is running on port %d with protocol %s\n",
		server.Port, server.Protocol)
}
