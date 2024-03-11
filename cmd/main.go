package main

import (
	"fmt"
	"github.com/cyunrei/portbridge/cmd/options"
	"github.com/cyunrei/portbridge/cmd/rules"
	"github.com/cyunrei/portbridge/pkg/forward"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var version string

func main() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})

	rs := parseOptionsToRules()
	done := make(chan struct{})
	signals := make(chan os.Signal, 1)
	errorOccurred := make(chan struct{})
	errorCount := int64(0)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	for _, r := range rs {
		r := r
		go func() {
			err := startForwarding(r)
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
		log.Infof("Received signal %v. Shutting down...\n", sig)
		close(done)
	}
}

func parseOptionsToRules() []rules.Rule {
	var opts options.Options
	var rs []rules.Rule
	parser := flags.NewParser(&opts, flags.None)
	_, err := parser.ParseArgs(os.Args)

	switch {
	case opts.RuleFile != "":
		rs, parseRulesErr := rules.ParseRulesFromFile(opts.RuleFile)
		if parseRulesErr != nil {
			log.Fatalf("Parse rules from file: %s", parseRulesErr)
		}
		if len(rs) == 0 {
			log.Fatalf("No rules found in the file '%s'. Please provide at least one rule", opts.RuleFile)
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
	case opts.RuleFile == "":
		rs = append(rs, rules.ParseRuleFromOptions(opts))
	}
	return rs
}

func startForwarding(r rules.Rule) error {
	fc := forward.NewForwardingConfig().WithSourceAddr(r.SourceAddr).
		WithDestinationAddr(r.DestinationAddr).WithProtocol(r.Protocol)
	switch r.Protocol {
	case "tcp":
		fc.WithDataForwarder(forward.NewTCPDataForwarder().SetBandwidthLimit(r.BandwidthLimit))
	case "udp":
		fc.WithDataForwarder(forward.NewUDPDataForwarder().SetBandwidthLimit(r.BandwidthLimit).
			SetDeadlineSecond(r.UDPTimeoutSecond).SetBufferSize(r.UDPBufferSize))
	}
	return fc.StartPortForwarding()
}
