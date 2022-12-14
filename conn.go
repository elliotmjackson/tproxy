package tproxy

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/elliotmjackson/tproxy/internal/display"
	"github.com/elliotmjackson/tproxy/internal/protocol"
	"github.com/elliotmjackson/tproxy/internal/writer"
	"github.com/fatih/color"
)

const (
	useOfClosedConn = "use of closed network connection"
)

type PairedConnection struct {
	id       int
	cliConn  net.Conn
	svrConn  net.Conn
	once     sync.Once
	stopChan chan struct{}
	protocol string
	quiet    bool
	delay    time.Duration
}

func NewPairedConnection(
	id int,
	cliConn net.Conn,
	protocol string,
	quiet bool,
	delay time.Duration,
) *PairedConnection {
	return &PairedConnection{
		id:       id,
		cliConn:  cliConn,
		stopChan: make(chan struct{}),
		protocol: protocol,
		quiet:    quiet,
		delay:    delay,
	}
}

func (p *Proxy) StartListener() error {
	conn, err := net.Listen("tcp", p.Local.Addr())
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	defer conn.Close()

	display.PrintfWithTime("Listening on %s...\n", conn.Addr().String())

	var connIndex int
	for {
		cliConn, err := conn.Accept()
		if err != nil {
			return fmt.Errorf("server: accept: %w", err)
		}

		connIndex++
		display.PrintlnWithTime(color.HiGreenString("[%d] Accepted from: %s",
			connIndex, cliConn.RemoteAddr()))

		pconn := NewPairedConnection(connIndex, cliConn, p.protocol, p.quiet, p.delay)
		go pconn.process(p.Remote)
	}
}

func (c *PairedConnection) copyData(dst io.Writer, src io.Reader, tag string) {
	_, e := io.Copy(dst, src)
	if e != nil && e != io.EOF {
		netOpError, ok := e.(*net.OpError)
		if ok && netOpError.Err.Error() != useOfClosedConn {
			reason := netOpError.Unwrap().Error()
			display.PrintlnWithTime(color.HiRedString("[%d] %s error, %s", c.id, tag, reason))
		}
	}
}

func (c *PairedConnection) handleClientMessage() {
	// client closed also trigger server close.
	defer c.stop()

	r, w := io.Pipe()
	tee := io.MultiWriter(c.svrConn, w)
	go protocol.CreateInterop(c.protocol).Dump(r, protocol.ClientSide, c.id, c.quiet)
	c.copyData(tee, c.cliConn, protocol.ClientSide)
}

func (c *PairedConnection) handleServerMessage() {
	// server closed also trigger client close.
	defer c.stop()

	r, w := io.Pipe()
	tee := io.MultiWriter(writer.NewDelayedWriter(c.cliConn, c.delay, c.stopChan), w)
	go protocol.CreateInterop(c.protocol).Dump(r, protocol.ServerSide, c.id, c.quiet)
	c.copyData(tee, c.svrConn, protocol.ServerSide)
}

func (c *PairedConnection) process(remote string) {
	defer c.stop()

	conn, err := net.Dial("tcp", remote)
	if err != nil {
		display.PrintlnWithTime(color.HiRedString("[x][%d] Couldn't connect to server: %v", c.id, err))
		return
	}

	display.PrintlnWithTime(color.HiGreenString("[%d] Connected to server: %s", c.id, conn.RemoteAddr()))

	c.svrConn = conn
	go c.handleServerMessage()

	c.handleClientMessage()
}

func (c *PairedConnection) stop() {
	c.once.Do(func() {
		close(c.stopChan)
		if c.cliConn != nil {
			display.PrintlnWithTime(color.HiBlueString("[%d] Client connection closed", c.id))
			c.cliConn.Close()
		}
		if c.svrConn != nil {
			display.PrintlnWithTime(color.HiBlueString("[%d] Server connection closed", c.id))
			c.svrConn.Close()
		}
	})
}
