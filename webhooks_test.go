package openroam

import (
	"context"
	"errors"
	"testing"
	"time"
)

const (
	vSecret = "whsec_test_secret"
	vTS     = "1770000000"
	vID     = "3db22f3e-8a20-4e64-9f89-5ab49ef39f4d"
	vSig    = "bb4f9bc4366d7f728b5e08c0f7f40749dc487166c30fb6c307bc374c77b8e320"
	vBody   = `{"eventId":"3db22f3e-8a20-4e64-9f89-5ab49ef39f4d","eventType":"partner.balance.low"}`
)

var vNow = time.Unix(1770000100, 0)

func TestAValidRequestVerifies(t *testing.T) {
	if err := VerifyWebhook(vSecret, vTS, vID, vSig, []byte(vBody), 0, vNow); err != nil {
		t.Fatal(err)
	}
}

func TestARefusedRequest(t *testing.T) {
	cases := map[string]error{
		"changed body": VerifyWebhook(vSecret, vTS, vID, vSig, []byte(vBody+" "), 0, vNow),
		"wrong secret": VerifyWebhook("other", vTS, vID, vSig, []byte(vBody), 0, vNow),
		"too old":      VerifyWebhook(vSecret, vTS, vID, vSig, []byte(vBody), 0, time.Unix(1770001000, 0)),
		"bad time":     VerifyWebhook(vSecret, "x", vID, vSig, []byte(vBody), 0, vNow),
	}
	for name, err := range cases {
		if !errors.Is(err, ErrWebhookSignature) {
			t.Errorf("%s: got %v", name, err)
		}
	}
}

func TestTheReadmeSetupCompilesAndBuildsARequest(t *testing.T) {
	cfg := NewConfiguration()
	cfg.Servers = ServerConfigurations{{URL: SandboxURL}}
	client := NewAPIClient(cfg)
	ctx := context.WithValue(context.Background(), ContextAPIKeys, map[string]APIKey{"PartnerApiKey": {Key: "key"}})
	_ = client.EsimsAPI.EsimsSearch(ctx).CountryCode("JP")
}
