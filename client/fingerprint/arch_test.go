package fingerprint

import (
	"testing"

	"github.com/actiontech/dtle/client/config"
	"github.com/actiontech/dtle/helper/testlog"
	"github.com/actiontech/dtle/nomad/structs"
)

func TestArchFingerprint(t *testing.T) {
	f := NewArchFingerprint(testlog.HCLogger(t))
	node := &structs.Node{
		Attributes: make(map[string]string),
	}

	request := &FingerprintRequest{Config: &config.Config{}, Node: node}
	var response FingerprintResponse
	err := f.Fingerprint(request, &response)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !response.Detected {
		t.Fatalf("expected response to be applicable")
	}

	assertNodeAttributeContains(t, response.Attributes, "cpu.arch")
}
