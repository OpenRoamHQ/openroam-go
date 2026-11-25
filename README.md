<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="brand/OpenRoam_Lockup_Color-On-Dark.svg">
    <img alt="OpenRoam" src="brand/OpenRoam_Lockup_Ink.svg" width="320">
  </picture>
</p>

# OpenRoam SDK for Go

Official client for the OpenRoam Partner API: eSIMs, airtime, utility payments, gift cards, hotels and flights.

## Install

```bash
go get github.com/OpenRoamHQ/openroam-go@v0.1.0-rc.7
```

> This is a **release candidate**: install it by its exact version, as above; a plain install does not pick it up. Before 1.0 the API can change between minor versions; every change is listed in the [changelog](https://openroamhq.github.io/openroam-sdks/changelog/).

## Environments

Partners use **Sandbox** (`https://api.sandbox.openroam.com`) or **Production** (`https://api.openroam.com`).
Default is Production. Pass `openroam.WithEnvironment(openroam.Sandbox)` for Sandbox. `WithBaseURL` is an unnamed
override only — not a third environment. The SDK sends your key as `X-Partner-Key`; a key works only on the environment
that issued it.

## Quick start

```go
import openroam "github.com/OpenRoamHQ/openroam-go"

client, err := openroam.New(os.Getenv("OPENROAM_API_KEY")) // or openroam.New("") with OPENROAM_API_KEY set
if err != nil {
	return err
}

found, err := client.Esims.Search(ctx, &openroam.EsimsSearchQuery{CountryCode: openroam.PtrString("JP")})

body := openroam.NewEsimOrderRequest(items) // items: the offers to buy, []openroam.EsimOrderRequestItem
body.SetCurrency("USD")
quote, err := client.Esims.Quote(ctx, body)
order, err := client.Esims.Order(ctx, body)

usage, err := client.Esims.GetUsage(ctx, iccid)
hotels, err := client.Hotels.Search(ctx, openroam.NewHotelSearchRequest( /* ... */ ))
```

Every product is a field of the client: `Esims`, `Airtime`, `Utilities`, `GiftCards`, `Hotels`, `Flights`, plus `Orders`, `OrderItems`,
`Wallets`, `Account`, `BrowserSessions`, `Checkouts`, `Currencies`, `Countries` and `MobileNetworks`. Buying always reads `Search`, `Quote`, `Order`; `Get` fetches one thing you own,
`List...` returns many, everything else starts with a verb. Every method takes `ctx` first, then the path ids in path order, then the
request body (a pointer to its model, passed as it is) or the query parameters (a pointer to the method's `...Query` struct; `nil` sends
none, and a `nil` field is not sent). A client is safe for concurrent use: make one and share it.

`Checkouts.Confirm` and `Checkouts.Fail` answer a checkout of a site that takes payment on your own page. Send them from your server
with a client made with that site's confirmation key (`orck_...`): a partner API key gets `403`. No signature headers are sent.

Options when you need them:

```go
client, err := openroam.New("key",
	openroam.WithEnvironment(openroam.Sandbox), // default openroam.Production
	openroam.WithMaxRetries(3),                 // default 0: no retries
	openroam.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	// openroam.WithBaseURL("https://...") // optional unnamed override (not a named environment)
)
```

## Idempotency

An idempotency key means: if the same request arrives twice, OpenRoam makes one order, not two. `Order` creates a key for you on every call.
To stay safe even if *your own* program crashes and asks again, pass a UUID generated once and persisted for that purchase, for example your order id if it is a UUID:

```go
order, err := client.Esims.Order(ctx, body, openroam.WithIdempotencyKey(myOrderID))
```

`openroam.NewIdempotencyKey()` makes a new one.

## Errors

A failed call returns an `*openroam.OpenRoamError` with the error body of the API: `Status`, `Code` (branch on this), `Name`, `Message`, `RequestID`
(keep it for support) and `Details`. When no answer arrived (network error, timeout, a cancelled `ctx`), `Status` is 0 and `errors.Unwrap`
gives the original error, so `errors.Is(err, context.DeadlineExceeded)` still works.

```go
order, err := client.Esims.Order(ctx, body)
var apiErr *openroam.OpenRoamError
if errors.As(err, &apiErr) && apiErr.Code == 11002 {
	// the price rose above your maxPrice
} else if err != nil {
	return err
}
```

A `202` whose body has code `12002` (the same order is still being accepted) is an error too, never an order.

## Lists and pages

Pagination depends on the method. `Esims.Search` accepts `Page` and `Limit` and returns `Items` and `Pagination`, while `Orders.List` returns `Orders` and `Pagination`.
`Esims.ListCatalog` accepts `Cursor` and `Limit` and returns `Items` and `NextCursor`; pass `NextCursor.Get()` as `Cursor` to fetch the next page.
Some lists, such as `Esims.ListCompatibleDevices`, return an unpaged slice.

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
    case item.EsimOrderItem != nil:
        fmt.Println(item.EsimOrderItem.Fulfillment) // eSIM details are kept
    case item.AirtimeOrderItem != nil:
        // ...
    }
}
```

An unknown `productType` is an error, never a silent gap.

## Retries

Retries are off by default. `WithMaxRetries` retries network errors, `429`, `502/503/504` and "still running" (`12002`) safely, with
backoff and `Retry-After`; `500` and other `4xx` never are. A `POST` is sent again only with its `Idempotency-Key`, the same key and the
same body, and `Order` always has one.

```go
client, err := openroam.New(apiKey, openroam.WithMaxRetries(3))
```

The full behaviour: [Retries and idempotency](https://openroamhq.github.io/openroam-sdks/retries/). `NewRetryingHTTPClient` (with `RetryOptions` for the delays
and an `OnRetry` callback for logging) and `NewIdempotencyKey` are also exported if you build your own requests.

## Good to know

- **Amounts** are `{ amount, currency }`. `price` is what OpenRoam charges you, `value` what your customer receives. Do not do arithmetic on amounts (JSON numbers are floating point in this language).
- API reference and partner portal: https://openroam.com. Support: support@openroam.com.

The client is generated from the OpenRoam Partner API description, so it always matches the API; `webhooks`, `retry` and the errors are hand-written. MIT licensed.
