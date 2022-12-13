package tproxy

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/fatih/color"
	"github.com/kevwan/tproxy/internal/config"
	"github.com/kevwan/tproxy/internal/display"
	"github.com/kevwan/tproxy/internal/protocol"
	"github.com/kevwan/tproxy/internal/writer"
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
	settings config.Settings
}

func NewPairedConnection(id int, cliConn net.Conn, settings config.Settings) *PairedConnection {
	return &PairedConnection{
		id:       id,
		cliConn:  cliConn,
		stopChan: make(chan struct{}),
		settings: settings,
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
	go protocol.CreateInterop(c.settings.Protocol).Dump(r, protocol.ClientSide, c.id, c.settings.Quiet)
	c.copyData(tee, c.cliConn, protocol.ClientSide)
}

func (c *PairedConnection) handleServerMessage() {
	// server closed also trigger client close.
	defer c.stop()

	r, w := io.Pipe()
	tee := io.MultiWriter(writer.NewDelayedWriter(c.cliConn, c.settings.Delay, c.stopChan), w)
	go protocol.CreateInterop(c.settings.Protocol).Dump(r, protocol.ServerSide, c.id, c.settings.Quiet)
	c.copyData(tee, c.svrConn, protocol.ServerSide)
}

func (c *PairedConnection) process() {
	defer c.stop()

	conn, err := net.Dial("tcp", c.settings.Remote)
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

func StartListener(settings config.Settings) error {
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", settings.LocalHost, settings.LocalPort))
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

		pconn := NewPairedConnection(connIndex, cliConn, settings)
		go pconn.process()
	}
}
