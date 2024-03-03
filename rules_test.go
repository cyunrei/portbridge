package main

import (
	"reflect"
	"testing"
)

func TestParseRulesFile(t *testing.T) {
	testCases := []struct {
		filePath string
		expected []Rule
	}{
		{
			filePath: "rules_example.json",
			expected: []Rule{
				{SourceAddr: "127.0.0.2:53", DestinationAddr: "[2606:4700:4700::1111]:53", Protocol: "udp"},
				{SourceAddr: "127.0.0.1:7777", DestinationAddr: "[::1]:7777", Protocol: "tcp"},
				{SourceAddr: ":8081", DestinationAddr: "127.0.0.1:8080", Protocol: "tcp", BandwidthLimit: 1024},
			},
		},
		{
			filePath: "rules_example.yaml",
			expected: []Rule{
				{SourceAddr: "127.0.0.2:53", DestinationAddr: "[2606:4700:4700::1111]:53", Protocol: "udp"},
				{SourceAddr: "127.0.0.1:7777", DestinationAddr: "[::1]:7777", Protocol: "tcp"},
				{SourceAddr: ":8081", DestinationAddr: "127.0.0.1:8080", Protocol: "tcp", BandwidthLimit: 1024},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.filePath, func(t *testing.T) {
			rules, err := parseRulesFromFile(testCase.filePath)
			if err != nil {
				t.Fatalf("Error parsing file %s: %v", testCase.filePath, err)
			}

			if !reflect.DeepEqual(rules, testCase.expected) {
				t.Errorf("Mismatched result for file %s.\nExpected: %v\nGot: %v", testCase.filePath, testCase.expected, rules)
			}
		})
	}
}
