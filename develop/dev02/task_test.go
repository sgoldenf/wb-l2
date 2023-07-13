package main

import (
	"testing"
)

func TestStringUnpack(t *testing.T) {
	tests := map[string]struct {
		input  string
		output string
		err    error
	}{
		"string with numbers": {
			input:  "a4bc2d5e",
			output: "aaaabccddddde",
			err:    nil,
		},
		"string without numbers": {
			input:  "abcd",
			output: "abcd",
			err:    nil,
		},
		"empty string": {
			input:  "",
			output: "",
			err:    nil,
		},
		"invalid string (only with number)": {
			input:  "45",
			output: "",
			err:    errorInvalidString,
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			unpacked, err := stringUnpack(v.input)
			if unpacked != v.output {
				t.Errorf("expected: %s, got: %s", v.output, unpacked)
			}
			if err != v.err {
				t.Errorf("expected: %v, got: %v", v.err, err)
			}
		})
	}
}
