package bcamqp

import (
	"testing"
)

func init() {
	// adaptorInstance = mock.Server{}
}

func TestConnInit(t *testing.T) {
	_, err := Open(BrokerOptions{}, func(err error) {
		t.Log(err)
	})

	if err != nil {
		t.Fatal(err)
	}
}
