package aws

import (
	"testing"
)

func TestNew(t *testing.T) {
	endpoint := ""
	broker, err := New(endpoint, "eu-west-1", "", "")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	token := broker.Publish("mytopic", "fuccccckkkk")
	if token.Wait() && token.Error() != nil {
		t.Fatalf("error: %s", token.Error())
	}
}
