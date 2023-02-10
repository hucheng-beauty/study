package builder_mode

import (
    "fmt"
    "testing"
)

func TestNewServerBuilder(t *testing.T) {
    builder := NewServerBuilder()
    server := builder.WithPort(8080).
        WithProtocol("http").
        Build()

    fmt.Printf("Server is running on port %d with protocol %s\n",
        server.Port, server.Protocol)
}
