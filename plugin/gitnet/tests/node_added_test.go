package tests

import (
	"testing"
)

var (
	endpoint  = "1.1.1.1.json"
	endpoint2 = "3.1.4.1.json"
	endpoints = []string{"1.1.1.1.json", "2.2.2.2.json", "3.3.3.3.json", "4.4.4.4.json", "5.5.5.5.json", "6.6.6.6.json"}
)

func TestEndpointDoesExist(t *testing.T) {
	if ok := !fileExists(endpoint, endpoints); ok {
		t.Errorf("mock endpoint [ %s ] should exist in the endpoint slice %s", endpoint, endpoints)
	}
}

func TestEndpointDoesNotExist(t *testing.T) {
	if ok := fileExists(endpoint2, endpoints); ok {
		t.Errorf("mock endpoint [ %s ] should not exist in the endpoint slice %s", endpoint2, endpoints)
	}
}

func fileExists(s string, endpoints []string) bool {
	for _, endpoint := range endpoints {
		if endpoint == s {
			return true
		}
	}
	return false
}
