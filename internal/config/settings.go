package config

import (
	"time"
)

type Settings struct {
	Remote    string
	LocalHost string
	LocalPort int
	Delay     time.Duration
	Protocol  string
	Stat      bool
	Quiet     bool
}

func SaveSettings(localHost string, localPort int, remote string, delay time.Duration,
	protocol string, stat, quiet bool,
) Settings {
	out := Settings{}
	if localHost != "" {
		out.LocalHost = localHost
	}
	if localPort != 0 {
		out.LocalPort = localPort
	}
	if remote != "" {
		out.Remote = remote
	}
	out.Delay = delay
	out.Protocol = protocol
	out.Stat = stat
	out.Quiet = quiet
	return out
}
