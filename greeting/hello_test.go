package greeting

import (
	"testing"
)

func TestHello(t *testing.T) {
	hello := Hello("AWSWesleyExternalClusterConnector")
	expected := "Hello World - AWSWesleyExternalClusterConnector"
	if hello != expected {
		t.Errorf("Expected %s, got %s", expected, hello)
	}
}
