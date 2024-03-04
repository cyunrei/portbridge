package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var version string
var forwarder TCPDataForwarder

type Options struct {
	SourceAddr      string `short:"s" long:"source" description:"Source address and port to bind locally" required:"true"`
	DestinationAddr string `short:"d" long:"destination" description:"Destination address and port to connect remotely" required:"true"`
	Protocol        string `short:"p" long:"protocol" description:"Specify the source protocol type" required:"true"`
	BandwidthLimit  int64  `short:"b" long:"bandwidth-limit" description:"TCP Bandwidth limit in KiB" default:"0"`
	RuleFile        string `short:"f" long:"rule-file" description:"Batch port forwarding file path"`
	Help            bool   `short:"h" long:"help" description:"Show help message"`
	Version         bool   `short:"v" long:"version" description:"Print the version number"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
	forwarder = &SimpleTCPDataForwarder{}
}

func main() {
	var opts Options
	var rules []Rule
	parser := flags.NewParser(&opts, flags.None)
	_, err := parser.ParseArgs(os.Args)
	switch {
	case opts.RuleFile != "":
		var parseRulesErr error
		rules, parseRulesErr = parseRulesFromFile(opts.RuleFile)
		if parseRulesErr != nil {
			log.Fatalf("Parse rules from file: %s", parseRulesErr)
		}
		if len(rules) == 0 {
			log.Fatalf("No rules found in the file '%s'. Please provide at least one rule", opts.RuleFile)
		}
		break
	case opts.Help:
		parser.WriteHelp(os.Stdout)
		return
	case opts.Version:
		fmt.Println("PortBridge Version:", version)
		return
	case err != nil:
		parser.WriteHelp(os.Stdout)
		log.Fatalf("Error: %s", err)
	case opts.RuleFile == "":
		rules = append(rules, parseRuleFromOptions(opts))
	}
	done := make(chan struct{})
	signals := make(chan os.Signal, 1)
	errorOccurred := make(chan struct{})
	errorCount := int64(0)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	for _, rule := range rules {
		rule := rule
		innerForwarder := forwarder
		if rule.BandwidthLimit > 0 {
			innerForwarder = &TrafficControlTCPDataForwarder{
				BandwidthLimit: rule.BandwidthLimit,
			}
			log.Infof("Forward TCP with bandwidth limit: %d KiB/s", opts.BandwidthLimit)
		}
		go func() {
			err := startPortForwarding(ForwardingConfig{
				SourceAddr:       rule.SourceAddr,
				DestinationAddr:  rule.DestinationAddr,
				Protocol:         rule.Protocol,
				TCPDataForwarder: innerForwarder,
			})
			if err != nil {
				log.Errorf("Error: %s", err)
				atomic.AddInt64(&errorCount, 1)
				if errorCount == int64(len(rules)) {
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
