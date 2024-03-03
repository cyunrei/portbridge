package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Rule struct {
	SourceAddr      string `json:"source_addr" yaml:"source_addr"`
	DestinationAddr string `json:"destination_addr" yaml:"destination_addr"`
	Protocol        string `json:"protocol" yaml:"protocol"`
	BandwidthLimit  int64  `json:"bandwidth_limit" yaml:"bandwidth_limit"`
}

func parseRulesFromFile(filePath string) ([]Rule, error) {
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
		return nil, fmt.Errorf("unsupported file format: %s\n", ext)
	}

	return rules, nil
}

func parseRuleFromOptions(opts Options) Rule {
	return Rule{
		SourceAddr:      opts.SourceAddr,
		DestinationAddr: opts.DestinationAddr,
		Protocol:        opts.Protocol,
		BandwidthLimit:  opts.BandwidthLimit,
	}
}
