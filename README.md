# OpenRoam SDK for Go

Official client for the OpenRoam Partner API: eSIMs, airtime, utility payments, gift cards, hotels and flights.

## Install

```bash
go get github.com/OpenRoamHQ/openroam-go@v0.1.0-rc.2
```

> This is a **release candidate**: install it by its exact version, as above; a plain install does not pick it up. Before 1.0 the API can change between minor versions; every change is listed in the [changelog](https://github.com/OpenRoamHQ/openroam-sdks/blob/main/CHANGELOG.md).

## Quick start

```go
import (
	"context"
	"os"

	openroam "github.com/OpenRoamHQ/openroam-go"
)

cfg := openroam.NewConfiguration()
cfg.Servers = openroam.ServerConfigurations{{URL: openroam.SandboxURL}} // or openroam.ProductionURL
client := openroam.NewAPIClient(cfg)

ctx := context.WithValue(context.Background(), openroam.ContextAPIKeys,
	map[string]openroam.APIKey{"PartnerApiKey": {Key: os.Getenv("OPENROAM_API_KEY")}})

result, _, err := client.EsimsAPI.EsimsSearch(ctx).CountryCode("JP").Execute()
```

Every product has its own service (`AirtimeAPI`, `EsimsAPI`, `GiftCardsAPI`, `HotelsAPI`, `FlightsAPI`, `UtilityAPI`, `OrdersAPI`, `WalletsAPI`...) and each follows search, quote, order.

## Webhooks

```go
err := openroam.VerifyWebhook(secret,
	r.Header.Get("X-OpenRoam-Timestamp"), r.Header.Get("X-OpenRoam-Event-Id"),
	r.Header.Get("X-OpenRoam-Signature"),
	rawBody, // the body exactly as received, before JSON parsing
	0, time.Time{}) // defaults: 5 minutes tolerance, now
if errors.Is(err, openroam.ErrWebhookSignature) { /* reject */ }
```

## Order items and webhook items

An order item (and a webhook item, and an item of an order request) is a union with one variant per product. Read it through the pointer for its product; exactly one is set, chosen by `productType`:

```go
for _, item := range order.Items {
    switch {
    case item.PartnerOrderItemEsim != nil:
        fmt.Println(item.PartnerOrderItemEsim.Fulfillment) // eSIM details are kept
    case item.PartnerOrderItemAirtime != nil:
        // ...
    }
}
```

An unknown `productType` is an error, never a silent gap.

## Good to know

- **Orders need an `Idempotency-Key`** (a new UUID per purchase). After a timeout or a `5xx`, retry with the same key and the same body: it never charges twice.
- **Amounts** are `{ amount, currency }`. `price` is what OpenRoam charges you, `value` what your customer receives.
- **Errors** carry `code`, `name`, `message`, `requestId` and `details`; quote `requestId` when you contact support.
- Full API documentation: https://openroam.com (partner portal). Support: support@openroam.com.

This SDK is generated from the OpenRoam Partner API description; `webhooks` is hand-written. MIT licensed.
