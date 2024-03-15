package rules

import (
	"encoding/json"
	"fmt"
	"github.com/cyunrei/portbridge/cmd/options"
	"github.com/cyunrei/portbridge/pkg/forwarder"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type Rule struct {
	SourceAddr      string `json:"source_addr" yaml:"source_addr"`
	DestinationAddr string `json:"destination_addr" yaml:"destination_addr"`
	Protocol        string `json:"protocol" yaml:"protocol"`
	BandwidthLimit  uint64 `json:"bandwidth_limit" yaml:"bandwidth_limit"`
	UDPBufferSize   uint64 `json:"udp_buffer_size" yaml:"udp_buffer_size"`
	UDPTimeout      uint64 `json:"udp_timeout" yaml:"udp_timeout"`
}

func NewRules() []Rule {
	return []Rule{
		{
			UDPBufferSize: forwarder.DefaultUDPBufferSize,
			UDPTimeout:    forwarder.DefaultUDPDeadlineTime,
		},
	}
}

func ParseFromFile(filePath string) ([]Rule, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var rules []Rule

	switch ext := filepath.Ext(filePath); ext {
	case ".json":
		if err := json.Unmarshal(file, &rules); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(file, &rules); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	return applyDefault(rules), nil
}

func ParseFromOptions(opts options.Options) Rule {
	return Rule{
		SourceAddr:      opts.SourceAddr,
		DestinationAddr: opts.DestinationAddr,
		Protocol:        opts.Protocol,
		BandwidthLimit:  opts.BandwidthLimit,
		UDPBufferSize:   opts.UDPBufferSize,
		UDPTimeout:      opts.UDPTimeout,
	}
}

func GenerateEmptyFile(filePath string) error {
	format := filepath.Ext(filePath)
	format = strings.TrimPrefix(format, ".")
	emptyRules := NewRules()

	var data []byte
	var err error
	switch format {
	case "yaml", "yml":
		data, err = yaml.Marshal(emptyRules)
	case "json":
		data, err = json.MarshalIndent(emptyRules, "", "    ")
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

func applyDefault(rules []Rule) []Rule {
	for i := range rules {
		if rules[i].UDPBufferSize == 0 {
			rules[i].UDPBufferSize = forwarder.DefaultUDPBufferSize
		}
		if rules[i].UDPTimeout == 0 {
			rules[i].UDPTimeout = forwarder.DefaultUDPDeadlineTime
		}
	}
	return rules
}
