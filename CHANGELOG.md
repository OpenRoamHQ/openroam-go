# Changelog

All OpenRoam SDKs share one version and one changelog. Entries are written by the release process from the
changes in the Partner API description.

## 0.1.0-rc.5

Fifth release candidate. The API, its routes, its JSON and the client methods are the same as 0.1.0-rc.4; the type names
change (below). Use the exact version (a release candidate is not picked up by a plain install); report problems to
support@openroam.com.

**Breaking: type names.** The request and response types now have names that say what they are, in every SDK. Nothing on
the wire changes: the API, its routes and its JSON are the same, and so are the client methods. Code that names a type (an
import, a type annotation, `new SearchHotelQuery(...)`) needs the new name.

- Backend words are gone from the names: `B2B`, `Partner`, `Command`, `Query`, `Api`, `Model`, `Cart`, and verbs such as `Get`.
- One type per product and concept, product first: `HotelOffer`, `HotelQuoteItem`, `HotelOrderItem`, `HotelFulfillment`.
- Request bodies are `<Subject><Action>Request`; a response is named for what it is, without `Response`.
- Types that had two names for the same shape are now one: `Order`, `Currency`, `EsimAddon`, `FlightAircraft`,
  `FlightSelectedServiceInput`, `ManualAmendmentResult`.

| Before | Now |
|---|---|
| `B2BOrderResponse`, `B2BOrderAcceptedResponse` | `Order` |
| `AirtimeCreateOrder200Response` (and the five other products) | `AirtimeOrder` (`EsimOrder`, `FlightOrder`...) |
| `B2BOrderListResponse`, `B2BOrderSummaryResponse` | `OrderList`, `OrderSummary` |
| `PartnerOrderItem`, `PartnerOrderItemHotel` | `OrderItem`, `HotelOrderItem` |
| `ItemHotelFulfillment` | `HotelFulfillment` |
| `CreateB2BOrderCommand`, `CreateB2BOrderItemCommand` | `OrderRequest`, `OrderRequestItem` |
| `CreateB2BOrderItemCommandHotel` | `HotelMixedOrderRequestItem` |
| `B2BMixedQuoteResponse`, `B2BMixedQuoteItem`, `B2BMixedQuoteItemHotel` | `OrderQuote`, `OrderQuoteItem`, `HotelMixedQuoteItem` |
| `HotelPartnerQuoteResponse`, `HotelPartnerQuoteItem`; `AirtimeQuoteResponse` | `HotelQuote`, `HotelQuoteItem`; `AirtimeQuote` |
| `AirtimeSearchResponse`, `HotelPropertySearchResponse`, `FlightSearchApiResponse` | `AirtimeSearchResult`, `HotelPropertySearchResult`, `FlightSearchResult` |
| `SearchHotelQuery`, `SearchB2BFlightsQuery` | `HotelSearchRequest`, `FlightSearchRequest` |
| `HotelCartInputs` (and the other products) | `HotelInputs` |
| `FlightFareOfferApi` and every other `…Api` flight type | `FlightOffer`, and the name without `Api` |
| `FlightPartnerSlice`, `FlightPartnerSegment` (as booked) | `FlightBookedSlice`, `FlightBookedSegment` |
| `AirtimeCatalogResponse` (and the other products) | `AirtimeCatalogPage` |
| `GiftCardItem`, `AirtimeOperatorItem` | `GiftCard`, `AirtimeOperator` |
| `B2BOrderItemAmendment`, `QuoteB2BAmendmentQuery`, `ConfirmB2BAmendmentRequest` | `OrderItemAmendment`, `AmendmentQuoteRequest`, `AmendmentConfirmRequest` |
| `CreateB2BManualAmendmentCommand` | `ManualAmendmentCreateRequest` |
| `RefundOrderItemRequest`, `UpdateB2BOrderItemRequest` | `OrderItemCancelRequest`, `OrderItemUpdateRequest` |
| `B2BEsimProfile`, `B2BEsimUsageResponse`, `B2BEsimUsageHistoryResponse` | `EsimProfile`, `EsimCurrentUsage`, `EsimUsageReadingList` |
| `GetWalletItem`, `GetWalletWithStatistics` | `WalletSummary`, `Wallet` |
| `CreateB2BWalletCommand`, `TransferB2BWalletsCommand` | `WalletCreateRequest`, `WalletTransferRequest` |
| `GetB2BAllCurrenciesCommandResponse` | `CurrencyList` |
| `PartnerErrorResponse`, `PartnerErrorDetail` | `ErrorResponse`, `ErrorDetail` |
| `B2BOrderItemRefundedWebhook`, `B2BPartnerBalanceLowWebhook` | `OrderItemRefundedWebhook`, `BalanceLowWebhook` |

## 0.1.0-rc.4

Fourth release candidate. Every SDK now has a single client for all products, and the API changes in a few places that break calls
made with 0.1.0-rc.2 (listed below). 0.1.0-rc.3 carried the same content but was never published: its GitHub release could not be
created, so there is no package under that version; use 0.1.0-rc.4. Use the exact version (a release candidate is not picked up by
a plain install); report problems to support@openroam.com.

**New**

- One client per SDK, with a method per route, named the same way in all eight languages: buying is always `search`, `quote`,
  `order`; `get` fetches one thing you own; `list...` returns many. TypeScript, PHP and Java: `new OpenRoam("key")`; Python:
  `OpenRoam("key")`; Go: `openroam.New("key")`; .NET: `new OpenRoamClient("key")`; Dart: `OpenRoam('key')`; Ruby:
  `OpenRoam::Client.new("key")`. Each README shows the quick start. In every language except TypeScript the generated API
  classes of 0.1.0-rc.2 are still in the package.
- Opt-in retries for network errors, 429, 502/503/504 and the in-progress code 12002, with backoff that honours `Retry-After`.
  Only reads and requests that carry an `Idempotency-Key` are repeated, so a purchase is never sent twice by mistake. Off by
  default; each SDK has a helper that makes an `Idempotency-Key`.
- `GET /amendments` lists the amendments across all your orders (`orderItems.listPartnerAmendments`).
- `GET /esims/{iccid}/profile` returns an eSIM's profile (`esims.getProfile`).
- Order items: every item in an order response now carries its `orderId` and the actions `cancellation` and `update`.
  `orderItems` has `get`, `update` and `cancel` for the items of every product, and hotel manual amendments are under
  `orderItems.manualAmendments` (`create`, `cancel`, `approve`, `decline`).
- Hotels: special requests (`specialRequestCodes` on an order, `specialRequestOptions` and `allowsSpecialRequests` on a rate),
  `guestNameRules`, `minimumSellingPrice`, a `cancellationTypes` search filter, and `replacements` in the catalog files.

**Breaking, compared with 0.1.0-rc.2**

- TypeScript no longer exports the generated API classes (`EsimsApi`, `HotelsApi`...) or `Configuration` from `@openroam/sdk`.
  Use `new OpenRoam({ apiKey, environment: 'sandbox' })`.
- Six routes are gone. Use the route that replaces each (all of them are new in this candidate):
  - `GET /esims/{iccid}` is now `GET /esims/{iccid}/profile`.
  - `POST /order-items/{id}/refund` is now `POST /order-items/{id}/cancel`; `actions.refund` on an item is now
    `actions.cancellation`.
  - `GET /order-items/{id}/cancellation-penalty` is now `GET /order-items/{id}/cancellation-preview`.
  - `POST /order-items/{id}/amendments/request` and `.../amendments/commit` are now `.../amendments/quote` and
    `.../amendments/confirm`. Confirming an amendment pays from the wallet balance only (`paymentMethod` takes `balance`).
  - `POST /orders/{orderId}/fulfillments/retry` is now `POST /orders/{orderId}/fulfillment/retry`.
- Six error names changed; their codes did not.

## 0.1.0-rc.2

Second release candidate. The API is the same as 0.1.0-rc.1; this candidate fixes how the SDKs are published. 0.1.0-rc.1 reached
PyPI, NuGet, pub.dev, Go and PHP but not npm (a workflow error) or Maven Central (its bundle was rejected: the Java SDK no longer
ships a tests jar and the publishing plugin is newer). The NuGet package now carries a readme, copyright and release notes.
Use the exact version (a release candidate is not picked up by a plain install); report problems to support@openroam.com.

## 0.1.0-rc.1

First release candidate of the OpenRoam Partner API SDKs for TypeScript/JavaScript, Python, Go, PHP, Java, .NET, Dart and Ruby,
generated from the Partner API description. A release candidate is not picked up by a plain install: use the exact version (see each
README). The stable 0.1.0 follows once the candidate has been tried; report problems to support@openroam.com.
