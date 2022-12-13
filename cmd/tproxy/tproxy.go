package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kevwan/tproxy"
	"github.com/kevwan/tproxy/internal/config"
)

func main() {
	var (
		localPort = flag.Int("p", 0, "Local port to listen on, default to pick a random port")
		localHost = flag.String("l", "localhost", "Local address to listen on")
		remote    = flag.String("r", "", "Remote address (host:port) to connect")
		delay     = flag.Duration("d", 0, "the delay to relay packets")
		protocol  = flag.String("t", "", "The type of protocol, currently support grpc")
		stat      = flag.Bool("s", false, "Enable statistics")
		quiet     = flag.Bool("q", false,
			"Quiet mode, only prints connection open/close and stats, default false")
	)

	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	flag.Parse()
	settings := config.SaveSettings(*localHost, *localPort, *remote, *delay, *protocol, *stat, *quiet)

	if len(settings.Remote) == 0 {
		fmt.Fprintln(os.Stderr, color.HiRedString("[x] Remote target required"))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if err := tproxy.StartListener(settings); err != nil {
		fmt.Fprintln(os.Stderr, color.HiRedString("[x] Failed to start listener: %v", err))
		os.Exit(1)
	}
}
