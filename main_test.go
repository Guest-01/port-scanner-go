package main

import (
	"reflect"
	"testing"
)

func TestParsePorts(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []int
	}{
		{name: "single port", input: "80", expected: []int{80}},
		{name: "multiple ports", input: "80,443,8080", expected: []int{80, 443, 8080}},
		{name: "invalid port", input: "80,abc", expected: nil},
		{name: "empty input", input: "", expected: nil},
		{name: "just commas", input: ",,", expected: nil},
		{name: "not in port range", input: "99999", expected: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := parsePorts(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("input: %v expected: %v result: %v", tc.input, tc.expected, result)
			}
		})
	}
}
