package writer

import (
	"errors"
	"io"
	"time"
)

type delayedWriter struct {
	writer   io.Writer
	delay    time.Duration
	stopChan <-chan struct{}
}

func NewDelayedWriter(writer io.Writer, delay time.Duration, stopChan <-chan struct{}) delayedWriter {
	return delayedWriter{
		writer:   writer,
		delay:    delay,
		stopChan: stopChan,
	}
}

func (w delayedWriter) Write(p []byte) (int, error) {
	if w.delay == 0 {
		return w.writer.Write(p)
	}

	timer := time.NewTimer(w.delay)
	defer timer.Stop()

	select {
	case <-timer.C:
		return w.writer.Write(p)
	case <-w.stopChan:
		return 0, errors.New("client canceled")
	}
}
