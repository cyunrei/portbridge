package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var version string

type Options struct {
	SourceAddr      string `short:"s" long:"source" description:"Source address and port to bind locally" required:"true"`
	DestinationAddr string `short:"d" long:"destination" description:"Destination address and port to connect remotely" required:"true"`
	Protocol        string `short:"p" long:"protocol" description:"Specify the source protocol type" required:"true"`
	Help            bool   `short:"h" long:"help" description:"Show this help message"`
	Version         bool   `short:"v" long:"version" description:"Print the version number"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
}

func main() {
	var opts Options
	parser := flags.NewParser(&opts, flags.None)
	_, err := parser.ParseArgs(os.Args)
	switch {
	case opts.Help:
		parser.WriteHelp(os.Stdout)
		return
	case opts.Version:
		fmt.Println("PortBridge Version:", version)
		return
	case err != nil:
		fmt.Println("Error:", err)
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	done := make(chan struct{})
	errorOccurred := make(chan struct{})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := startPortForwarding(opts.SourceAddr, opts.DestinationAddr, opts.Protocol)
		if err != nil {
			log.Errorf("Error: %s\n", err)
			close(errorOccurred)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-errorOccurred:
		os.Exit(1)
	case sig := <-signals:
		log.Infof("Received signal %v. Shutting down...\n", sig)
		close(done)
	}
}
