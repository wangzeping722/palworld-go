//go:build linux

package main

import (
	"os"
	"os/signal"
	"syscall"
)

// ListenSignal 监听信号
func listenSignal() os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	return <-ch
}
