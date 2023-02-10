package builder_mode

type Option struct {
    Port     int
    Protocol string
}

type Server struct {
    Option
}

type ServerBuilder struct {
    option Option
}

func (b *ServerBuilder) WithPort(port int) *ServerBuilder {
    b.option.Port = port
    return b
}

func (b *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
    b.option.Protocol = protocol
    return b
}

func (b *ServerBuilder) Build() *Server {
    return &Server{b.option}
}

func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{}
}
