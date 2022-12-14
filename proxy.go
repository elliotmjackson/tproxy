package tproxy

import (
	"fmt"
	"time"
)

type Proxy struct {
	Local    Connection
	Remote   string
	protocol string
	quiet    bool
	delay    time.Duration
}

func NewProxy(localHost string, localPort int, remote string, protocol string, quiet bool, delay time.Duration) *Proxy {
	return &Proxy{
		Local: Connection{
			Host: localHost,
			Port: fmt.Sprintf("%d", localPort),
		},
		Remote:   remote,
		protocol: protocol,
		quiet:    quiet,
		delay:    delay,
	}
}

type Connection struct {
	Host string
	Port string
}

func (c Connection) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
