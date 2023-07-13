package main

import (
	"testing"
)

func TestGetTime(t *testing.T) {
	t.Run("valid server", func(t *testing.T) {
		_, err := getTime("0.beevik-ntp.pool.ntp.org")
		if err != nil {
			t.Errorf("expected: error = nil, got: %s", err.Error())
		}
	})
	t.Run("invalid server", func(t *testing.T) {
		_, err := getTime("0.beevik-ntp.pool.ntp.or")
		errString := "lookup 0.beevik-ntp.pool.ntp.or: no such host"
		if err.Error() != errString {
			t.Errorf("expected: %s, got: %s", errString, err.Error())
		}
	})
}
