// Hand-written, not generated (see overlay/): verify the signature of a webhook from OpenRoam.

package openroam

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

// The base URLs of the two environments.
const (
	ProductionURL = "https://b2b-api.openroam.com"
	SandboxURL    = "https://b2b-api-sandbox.openroam.com"
)

// ErrWebhookSignature means the request was not signed by OpenRoam with this secret, or is too old.
var ErrWebhookSignature = errors.New("openroam: invalid webhook signature")

// VerifyWebhook returns ErrWebhookSignature unless the request is validly signed and recent.
// timestamp, eventID and signature are the X-OpenRoam-Timestamp, X-OpenRoam-Event-Id and X-OpenRoam-Signature
// headers; rawBody is the request body exactly as received. A zero tolerance means 5 minutes; a zero now means now.
func VerifyWebhook(secret, timestamp, eventID, signature string, rawBody []byte, tolerance time.Duration, now time.Time) error {
	sentAt, err := strconv.ParseFloat(timestamp, 64)
	if err != nil {
		return ErrWebhookSignature
	}
	if tolerance == 0 {
		tolerance = 5 * time.Minute
	}
	if now.IsZero() {
		now = time.Now()
	}
	if math.Abs(float64(now.Unix())-sentAt) > tolerance.Seconds() {
		return ErrWebhookSignature
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "." + eventID + "."))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(strings.ToLower(signature))) {
		return ErrWebhookSignature
	}
	return nil
}
