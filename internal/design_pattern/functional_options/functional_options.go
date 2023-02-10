package functional_options

type Option struct {
	Port     int
	Protocol string
}

type Server struct {
	Option
}

type ServerOption func(*Option)

func WithPort(port int) ServerOption {
	return func(o *Option) {
		o.Port = port
	}
}

func WithProtocol(protocol string) ServerOption {
	return func(o *Option) {
		o.Protocol = protocol
	}
}

func NewServer(opts ...ServerOption) *Server {
	o := &Option{}
	for _, opt := range opts {
		opt(o)
	}
	return &Server{*o}
}
