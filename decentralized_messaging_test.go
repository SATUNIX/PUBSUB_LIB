package main

import (
	"testing"
)

func TestSomeFunction(t *testing.T) {
	result, err := someFunction()
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}

	if result != "" {
		t.Errorf("Expected an empty string, but got: %s", result)
	}
}