package main

import (
	"fmt"
	"github.com/cyunrei/portbridge/cmd/options"
	"github.com/cyunrei/portbridge/cmd/rules"
	"github.com/cyunrei/portbridge/pkg/forwarder"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var version string
var opts options.Options

func main() {

	rs := parseOptions()
	done := make(chan struct{})
	signals := make(chan os.Signal, 1)
	errorOccurred := make(chan struct{})
	errorCount := int64(0)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	for _, r := range rs {
		r := r
		go func() {
			err := startForwarder(r)
			if err != nil {
				log.Errorf("Error: %s", err)
				atomic.AddInt64(&errorCount, 1)
				if errorCount == int64(len(rs)) {
					close(errorOccurred)
				}
			}
		}()
	}
	select {
	case <-done:
	case <-errorOccurred:
		os.Exit(1)
	case sig := <-signals:
		log.Infof("Received signal %v. Shutting down...", sig)
		close(done)
	}
}

func parseOptions() []rules.Rule {
	var rs []rules.Rule
	parser := flags.NewNamedParser("portbridge", flags.None)
	parser.AddGroup("Options", "Application Options", &opts)
	_, err := parser.ParseArgs(os.Args)
	parser.Usage = "[OPTIONS]"
	initLog()
	switch {
	case opts.GenRulesFile:
		generateEmptyRulesFile()
		os.Exit(0)
	case opts.RulesFile != "":
		rs, parseRulesErr := rules.ParseFromFile(opts.RulesFile)
		if parseRulesErr != nil {
			log.Fatalf("Parse rules from file: %s", parseRulesErr)
		}
		if len(rs) == 0 {
			log.Fatalf("No rules found in the file '%s'. Please provide at least one rule", opts.RulesFile)
		}
		return rs
	case opts.Help:
		parser.WriteHelp(os.Stdout)
		os.Exit(0)
	case opts.Version:
		fmt.Println("PortBridge Version:", version)
		os.Exit(0)
	case err != nil:
		parser.WriteHelp(os.Stdout)
		log.Fatalf("Error: %s", err)
	case opts.RulesFile == "":
		rs = append(rs, rules.ParseFromOptions(opts))
	}
	return rs
}

func startForwarder(r rules.Rule) error {
	f := forwarder.NewForwarder().WithSourceAddr(r.SourceAddr).
		WithDestinationAddr(r.DestinationAddr).WithProtocol(r.Protocol)
	switch r.Protocol {
	case "tcp":
		df := forwarder.NewTCPDataForwarder().WithBandwidthLimit(r.BandwidthLimit)
		f.WithDataForwarder(df)
	case "udp":
		df := forwarder.NewUDPDataForwarder().WithBandwidthLimit(r.BandwidthLimit).
			WithDeadlineTime(r.UDPTimeout).WithBufferSize(r.UDPBufferSize)
		f.WithDataForwarder(df)
	}
	return f.Start()
}

func generateEmptyRulesFile() {
	err := rules.GenerateEmptyFile("example.yaml")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	err = rules.GenerateEmptyFile("example.json")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
	if opts.LogFile != "" {
		file, err := os.OpenFile(opts.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.Fatal("Failed to open log file: ", opts.LogFile, ", Error: ", err)
		}
	} else {
		log.SetOutput(os.Stdout)
	}
}
