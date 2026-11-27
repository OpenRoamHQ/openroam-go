# Changelog

All OpenRoam SDKs share one version and one changelog. Entries are written by the release process from the
changes in the Partner API description.

## 0.1.0-rc.7

**SDK hygiene.** Empty or whitespace-only webhook secrets are rejected in every language. C# generated
defaults no longer point at `localhost:3000`. OpenAPI env contrast that said Sandbox vs “Live” is rewritten to
“Production” on prepare (product-state `live` / `SITE_NOT_LIVE` unchanged). Partner READMEs warn against amount
arithmetic; the public docs site omits the Unreleased section.

**Breaking: an item's own fields are on the item, not under `details`.** On search items, quote items, mixed-quote items,
the catalog hotel and two search answers, the fields that were under `details` now sit beside the common keys, with no
alias: `{ id, name, offers, details: { destinationCode } }` is now `{ id, name, offers, destinationCode }`.

| Type | Now on the type itself |
| --- | --- |
| `AirtimeSearchResult` | `phoneLookup` |
| `AirtimeQuoteItem`, `AirtimeMixedQuoteItem` | `estimatedDeliveredValue` |
| `EsimDestination` | `destinationCode`, `region`, `coverage` |
| `GiftCard` | `isRedeemableGlobally`, `redemptionInstructions` |
| `Biller` | `serviceType`, `deliveryType` |
| `UtilityQuoteItem`, `UtilityMixedQuoteItem` | `deliveryType` |
| `Hotel`, `CatalogHotel` | `starRating`, `address`, `coordinates`, `imageUrls`, `amenityGroups`, `contact`, `displayGroup` |
| `HotelSearchResult` | `availableFilters` |
| `HotelQuoteItem`, `HotelMixedQuoteItem` | `expiresAt`, `availability`, `checkInDate`, `checkOutDate`, `requestedRooms`, `terms`, `rate` |
| `Itinerary` | `owner`, `slices`, `marketingCarriers`, `operatingCarriers`, `totalDurationMinutes`, `totalStops`, `totalEmissionsKg` |
| `FlightPlace` | `iataCode`, `cityName`, `cityIataCode`, `timezone`, `airports` |
| `FlightQuoteItem`, `FlightMixedQuoteItem` | `expiresAt`, `quotedAt`, `hasSeatMaps`, `terms`, `includedBaggage`, `availableServices`, `fareOptions` |

- Removed: the types `AirtimeSearchDetails`, `AirtimeQuoteItemDetails`, `EsimDestinationDetails`, `GiftCardDetails`,
  `BillerDetails`, `UtilityQuoteItemDetails`, `HotelDetails`, `HotelSearchDetails`, `HotelQuoteItemDetails`,
  `ItineraryDetails`, `FlightPlaceDetails` and `FlightQuoteItemDetails`. Constructors of the types above take the moved
  values in place of `details`.
- An enum declared on a removed type now belongs to the type that has the field. In .NET, `BillerDetails.ServiceTypeEnum`
  and `BillerDetails.DeliveryTypeEnum` are `Biller.ServiceTypeEnum` and `Biller.DeliveryTypeEnum`,
  `HotelQuoteItemDetails.AvailabilityEnum` is `HotelQuoteItem.AvailabilityEnum` or
  `HotelMixedQuoteItem.AvailabilityEnum`, and `UtilityQuoteItemDetails.DeliveryTypeEnum` is
  `UtilityQuoteItem.DeliveryTypeEnum` or `UtilityMixedQuoteItem.DeliveryTypeEnum`. Their numbers do not change.
- The hotel catalog files are now line format 3, with a hotel's facts on the line itself. `hotels.exportCatalog`
  (`GET /v1/hotels/catalog/files`) lists no files until the first build in this format.
- Unchanged: offers, flight fare options, purchased products, fulfillment results and webhooks keep `details`, and an
  error keeps its `details` list.

**Breaking: the amendment, manual amendment and direct update methods are gone.** Change an item with `orderItems.change`
and cancel it with `orderItems.cancel` (both below).

- `orderItems.quoteAmendment`, `confirmAmendment`, `listAmendments`, `getAmendment` and `listPartnerAmendments`
  (`/v1/order-items/{id}/amendments...` and `GET /v1/amendments`).
- `orderItems.manualAmendments` with `create`, `cancel`, `approve` and `decline`
  (`/v1/order-items/{id}/manual-amendments...`).
- `orderItems.update` (`POST /v1/order-items/{id}/update`).
- `orders.getCancellationPreview` (`GET /v1/orders/{id}/cancellation-preview`). Quote each item's cancellation instead.
- Their request and result types, `Amendment...`, `ManualAmendment...`, `OrderCancellationPreview...`,
  `OrderItemAmendment` and `OrderItemUpdateRequest`. `AmendmentQuoteOption` stays: it is an option of a change quote.

**New: `checkouts.confirm` and `checkouts.fail`, sent with the site's confirmation key.** For a site that takes payment
on your own page, your server confirms a checkout with `checkouts.confirm(id, { amount, paymentReference })`
(`POST /v1/checkouts/{id}/confirm`), or reports that the traveller's payment failed with `checkouts.fail(id)`
(`POST /v1/checkouts/{id}/fail`; it has no member, so in Go the method takes no body and sends `{}`). Send both with a client made with that site's confirmation
key (`orck_…`, issued in the dashboard); a partner API key is `403` / `10003` `KEY_NOT_ALLOWED`. Nothing is signed: the
`X-OpenRoam-Timestamp` and `X-OpenRoam-Signature` headers are gone from the confirmation and are not read. A failed
report makes the checkout `failed` with the new `failure.reason` `paymentFailed` (in .NET
`CheckoutFailure.ReasonEnum.PaymentFailed` is `12`; no other member changes its number).

**Breaking: an order item is cancelled in two steps on one route.** `orderItems.cancel(id, body)` now sends
`POST /v1/order-items/{id}/cancellation` with `OrderItemCancellationRequest`: `operation` `quote` returns what cancelling
now would keep and refund (`OrderItemCancellation.quote`, an `OrderItemCancellationPreview`), and `operation` `confirm`
with the `reason` and the quoted `refundAmount` cancels the item (`OrderItemCancellation.refund`). The routes
`POST /v1/order-items/{id}/cancel` and `GET /v1/order-items/{id}/cancellation-preview`, the method
`orderItems.getCancellationPreview` and the type `OrderItemCancelRequest` are gone, and
`OrderItemCancellationPreview.checkInDate` is removed. In .NET the reasons are `OrderItemCancellationRequest.ReasonEnum`
(`OrderItemCancelRequest.ReasonEnum` is gone).

**New: order item changes.** `orderItems.change(id, body)` (`POST /v1/order-items/{id}/changes`, `OrderItemChangeRequest`:
`request`, `confirm`, `decline` or `cancel`) returns `OrderItemChange`; `orderItems.getChange(id, changeId)`,
`orderItems.listChanges(id)` and `orderItems.listPartnerChanges()` (`GET /v1/changes`) read them. An item's
`actions.change` says whether a change can start. The `order-item.change.status` webhook (`OrderItemChangeStatusWebhook`)
replaces `order-item.manual-amendment.status`; the types `ManualAmendmentStatusWebhook`, `ManualAmendmentStatusWebhookData`
and `ManualAmendmentWallet` are gone.

**New fields.** An order placed from a checkout paid on your own page has `checkout` (`OrderCheckout`: the checkout `id`
and your `paymentReference`) and each item's `sellingPrice`; `order-item.refunded` carries both. A checkout shows the
`promoCode` it used (`CheckoutPromoCode`). An order item's `refund.wallet` is now `OrderItemRefundWallet`, which adds
`walletId`, and its `refund.methods` is new. `payment.card.last4` (`OrderPaymentCard`), `payment.wallet.walletId`,
`AmendmentQuoteOption.travellerRefund` and `travellerRefundCheckoutTax`, `GuestNameRules.guestNamesOptional` and
`HotelCatalogExport.linksExpireAt` are new.

**New: `flights.listPlaces`, the city and airport lookup.** `GET /v1/flights/places` finds the cities and airports to search
between, by name or IATA code (`q`, required), optionally only `category` `city` or `airport`, with `limit` (1 to 100, default
10). It returns `FlightPlaceList` (`items` of `FlightPlace`: `id`, `name`, `countryCode`, `category`, `logoUrl`, `iataCode`,
`cityName`, `cityIataCode`, the IANA `timezone` and, for a city, its `airports`). Pass a place's `id`, its three-letter
IATA code, as a slice's `origin` or `destination` in `flights.search`; a city stands for all of its airports. In .NET the new
enum is `FlightPlace.CategoryEnum` (`City` 1, `Airport` 2).

**Breaking: one tax block for every product.** Every offer, flight fare option, quote and mixed-quote item, order item
and checkout item now has `taxes` (`TaxSummary`): the taxes inside one unit of the `price` beside it, as `treatment`
(`included`, `addedOnTop`, `notStated` or `notApplicable`), `amount` (the complete tax, or null when it is not known;
never a partial figure and never a 0 that stands for "unknown") and `lines` (`TaxLine`, for display only; never add them
up). Airtime, utility, gift cards and eSIM state `notStated`.

- Removed: the types `Tax` and `TaxBreakdown` (a hotel rate's `details.taxes`, on offers and quoted rates) and
  `PriceBreakdown` (a flight's `details.priceBreakdown` and `passengerPricing[].priceBreakdown`). Read `taxes` instead.
  A tax a hotel collects at the property is a `taxes.lines` entry with `payableAt: onSite`, not a surcharge.
- A flight `PriceLine` no longer has the `taxes` category: the `fare` line includes the fare's taxes. In .NET,
  `PriceLine.CategoryEnum.Taxes` (2) is gone and the other members keep their numbers.
- New: `payment.checkoutTax` (`CheckoutTax`) on an order a traveller paid by card on one of your hosted sites: the part of
  `payment.amount` that is the tax OpenRoam added at checkout. It is already counted in `payment.amount`, so never add it
  again, and it is never a product's own tax.
- The API now also lists a `browserSession` bearer scheme, which a web page uses with a short-lived token. The SDKs are
  for your server and keep sending your API key.

**Breaking: every id has a prefix that says what it names.** Treat each id as an opaque string and pass it back
unchanged.

- Search items: an eSIM destination is `destn_…`, a gift card `card_…`, an airtime operator `operator_…` and a utility
  biller `biller_…`. `esims.getDestination`, `giftCards.getGiftCard`, `airtime.getOperator` and `utilities.getBiller` take
  these ids. An item id is never orderable: quote and order an offer's `publicId`.
- An airtime operator's `id` (`AirtimeOperator`, `AirtimeOperatorSummary`) is a string, not a UUID: `str` in Python (was
  `UUID`), `String` in Java (was `UUID`), `string` in .NET (was `Guid`; the constructor now refuses `null`). TypeScript,
  Go, PHP, Dart and Ruby already read it as a string.
- Offers: a hotel rate's `publicId` is `rate_…` (hotel quote returns a new `rate_…` id for the rate at its price), and a
  utility offer's `publicId` is an opaque `utility_…` (it was `utility_<biller>~<option>`). Other offer ids do not change.
- A flight itinerary's `id` is `itin_…`. It names the itinerary within one search answer only; quote and order an offer's
  `publicId`.

**New: hosted-site orders and account services.**

- `orders.resendEmail(id, { type, orderItemId })` sends an email of a hosted-site order to its traveller again: the
  `receipt`, or one item's `delivery` email. It needs an `Idempotency-Key`, which the SDK creates unless you pass yours. It
  returns `OrderEmailResendResult` (`orderId`, `type`, `orderItemId`, and `recipient`, masked).
- An order placed on one of your hosted sites carries `customer` (`OrderCustomer`: masked `email` and `phone`, `name`,
  `language`, `countryCode`). A payment's `status` can read `authorized`: the traveller's card is held for the order and is
  charged for what is delivered. In .NET, `OrderPayment.StatusEnum.Authorized` is `7`; no other member changes its number.
- `account.listServices()` lists every service with `enabled` (the account may sell it now) and `requestable` (the account
  can ask for it). A new `account` resource holds it.

**Breaking: orders, payments and refunds read as a developer expects.**

- An order's `status` no longer has `paid`: an order still in progress reads `processing`. Two statuses are new,
  `partiallyRefunded` and `refunded`. A payment's `status` spells `cancelled` (was `canceled`) and can also read
  `partiallyRefunded`.
- Every order and order summary has a `refunded` amount, and `effectivePaidAmount` is what you paid net of refunds. A refund,
  and the result of cancelling an item, names the wallet it went back to in `wallet` (`label`, and the amount `credited`).
- A rejected item is refunded without a request. The order-level `actions` and the two retry routes are gone, and with them
  `orders.retryPayment` and `orders.retryFulfillment` (and the types `OrderActions`, `OrderPaymentRetryResult` and
  `OrderFulfillmentRetryResult`).
- A wallet charge that fails no longer leaves an order behind: placing an order, or confirming an amendment, answers `402`
  with code `16000` (`PAYMENT_FAILED`), and `details[0]` names `wallet` with the reason `INSUFFICIENT_BALANCE` or
  `WALLET_BLOCKED`. Send the same request with the same `Idempotency-Key` after a top-up. Order routes no longer answer `202`.
- A request whose `Idempotency-Key` is still being processed answers `409` with code `12002` and `Retry-After` (was `202`); a
  finished one answers again with its first answer. The SDKs' retry layer already retries this `409`.
- An unsupported `currency` answers `400` with code `10001` at `currency`, reason `NOT_SUPPORTED` (was `404`).
- `orders.list` filters by `channel` (`api` or `dashboard`) instead of `salesChannel`.
- One spelling per value: an amendment refund reads `refunded` (was `completed`); the manual amendment webhook uses `quoted`,
  `completed` and `failed` (were `quoteReady`, `confirmed` and `rejected`); the refund webhook's `reason` is always present
  and is one of `requestedByCustomer`, `duplicate`, `fraudulent`, `fulfillmentFailed`. An action's `availability` no longer
  has `requiresDashboard`.

**Times: every instant is UTC, with the local clock beside it.** A flight's `departureAt`, `arrivalAt`, `arrivingAt` and
`departingAt` are UTC instants, and each now has a `...LocalTime` partner (`departureLocalTime`, `arrivalLocalTime`,
`arrivingLocalTime`, `departingLocalTime`) on segments, slices and stops, in search results and on booked flights: the
airport's clock, written `yyyy-MM-ddTHH:mm:ss` without a zone, for display only. The booking terms, the same shape for a
hotel stay and a flight, state the clock too: the cancellation terms gain `timezone` (the hotel's, or the first departure
airport's) and `freeUntilLocalTime`, and each step of the schedule `fromLocalTime`. A cancellation preview's penalties gain
`fromLocalTime` and `timezone`, and its `checkInDate` is a date without a time.

**Breaking: eSIMs.**

- One status pair everywhere: `esimStatus` (the eSIM's lifecycle) and `profileStatus` (the installed profile) replace
  `status` on the profile, the usage readings, the status history and the usage webhook. The SIM card on the profile has no
  `status` any more.
- No more sentinel numbers: an allowance or usage amount that is unlimited or not known is `null`, and unlimited is said by
  `isUnlimited` (`isDataUnlimited`, `isSmsUnlimited`, `isVoiceUnlimited` on usage). SMS and voice allowances have the same
  shape as data: `amount`, `unit`, `label`, `type`, `isUnlimited` (were `smsAmount`... and `callAmount`...). The profile's data
  allowance is an `EsimDataAllowance` (`EsimProfileData` is gone), and offer details and profiles lose their own
  `isUnlimited`. Throttle values are `null` when the plan has none.
- `expiredAt` is now `expiresAt`. A daily usage entry's `date` is a date, with `start` and `end` as date-times, and its
  `location` can be `null`. The usage and profile status webhooks send `start` and `end` as ISO 8601 date-times (were Unix
  seconds).
- `esims.sendSms` returns only when the message was sent: otherwise `422` with code `22000` (`ESIM_SMS_NOT_SUPPORTED`, do
  not retry), `502` with code `10008`, or `404`. Its result loses `success` and `isSupported`. eSIM search, quote, order and
  provisioning, and the quote and order of several products, declare the `404` they answer for an eSIM that is not yours.

**Other changes.**

- Hotel search has its own paging object, `HotelSearchPagination` (`currentCount`, `limit`, `hasMore`).
- A hotel's amendment capabilities say whether the dates (`dates.change`) and the holder (`holder.change`) may change.
- A flight's `pricePerPassenger` is `null` when the fare does not state one (it was the total split evenly).
- `wallets.update` changes only the fields you send, and its body no longer has `walletId`.
- Every webhook can carry `isTest: true`, for a signed test delivery; live events omit it.
- .NET: the new optional local time parameters are added at the end of their constructors (`Segment`, `Slice`, `Stop`, their
  booked forms and the cancellation types), so existing positional calls keep their meaning.
- .NET: enum members keep their numbers although the API lists some values in a new order: `WalletTransaction.TypeEnum`
  stays `Credit = 1` and `Debit = 2`, and every value an order status, payment status or action availability already had
  keeps its number. New values take new numbers (`PartiallyRefunded = 8` and `Refunded = 9` in the order statuses), and the
  numbers of removed values (`Paid = 2`) are not reused. A compiled comparison keeps matching the same value.
- Go: a union such as `OrderItem` reads hotel and flight items that carry a local time. The generated pattern check on
  nullable strings failed on every value, so it is not applied to them.

**Wallet transactions: debit takes money out, credit adds it.** The API had these the other way round. An order payment is
now a `debit` and a top-up or a refund is a `credit`, in a transaction's `type`, in the `type` filter and in a wallet's
`totalDebits` and `totalCredits`. The names and types are the same, so code compiles unchanged: if you read these values,
swap what you do with them. A transaction's `description` is now always present and is written by the API ("Wallet top-up",
"Order payment for order ..."), and a new `orderId` holds the order a payment or refund belongs to.

**Cancellation preview.** `checkInDate` can now be `null`: a flight has no check-in date, and neither has an item that was
never delivered. The preview of one item can answer `409` when the item cannot be cancelled now, with the reason, as the
cancel itself does; before, it could offer a refund the cancel then refused.

**Shorter routes.** Every route loses its `/api` and `/b2b` parts: `/api/v1/b2b/esims/search` is now `/v1/esims/search`, and
the same for every other route. SDK users change nothing: the client methods, their parameters and the types are the same,
and only the URL each method calls is new. Use this version once the API serves the new routes; from then on, older SDK
versions get `404` on every call. If you call the API without an SDK, change the paths in your code.

**New API hosts.** The SDKs now default to the new production and sandbox hosts:

| Environment | Before | Now |
|---|---|---|
| Production | `https://b2b-api.openroam.com` | `https://api.openroam.com` |
| Sandbox | `https://b2b-api-sandbox.openroam.com` | `https://api.sandbox.openroam.com` |

Code that uses the default host or chooses an environment by name (`environment: 'sandbox'`, `SandboxUrl`,
`OPENROAM_ENVIRONMENTS["sandbox"]` ...) gets the new host with this version. Code that passes an explicit base URL with one of
the old hosts should switch to the new one.

**Breaking: fields the API always returns are no longer optional.** The API description now marks every field an answer
always carries as required, in 253 response and webhook types, and every SDK types those fields as present. Nothing on the
wire changes: the routes, the JSON and the client methods are the same. Fields that can be missing or null (a `description`,
a `payment`, a `delivery`...) stay optional.

What you gain: `order.items`, `offer.price` or `fulfillment.details.codes` are there without a null check, an optional chain
or a pointer dereference, and your editor and type checker show which fields you can rely on.

What can break when you upgrade is code that builds these models itself, typically test fixtures and mocks, and code
written around the old optional types:

| Language | What changed | What to do |
|---|---|---|
| TypeScript | Always-returned properties lose their `?` (`publicId: string`, `details: HotelOfferDetails`). | An object literal typed as one of these types needs every required property. Code that reads them still compiles. |
| Python | The fields are no longer `Optional[...] = None`. | Building a model (`GiftCardOrder(...)`, `from_dict`) without one raises a pydantic `ValidationError`. Type checkers see the non-optional types. |
| Go | Required fields are values, not pointers (`Id *string` is now `Id string`, `Total *Money` is `Total Money`, an offer's `Value *OfferValue` is `Value OfferValue`). Their `HasX()` methods are gone, and `NewX()` takes the required fields as arguments (`NewGiftCardOrder(id, origin, status, total, ...)`). `GetX()`, `GetXOk()` and `SetX()` keep their signatures. | Drop dereferences (`*order.Id`) and `nil` checks on these fields, and pass the required fields to `NewX()`. |
| Java | No signature changes. The fields are annotated `@Nonnull` instead of `@Nullable`, and a model you serialize writes them even when they are null. | Nothing. |
| .NET | Enum properties of required fields lose their `?` (`StatusEnum? Status` is now `StatusEnum Status`; 79 properties). The constructors keep their parameters, but throw `ArgumentNullException` when a required object, list or string is null, so `new GiftCardOrder()` or `new Money()` with its defaults now throws. `Guid`, `DateTime`, number and `bool` properties were already non-nullable. | Replace `.HasValue` and `.Value` on those enums, and pass every required field to the constructor. |
| PHP | Getters and setters of required fields drop their `?` (`getId(): string`, `setId(string $id)`). | Do not pass `null` to these setters. `listInvalidProperties()` and `valid()` now name a missing required field. |
| Dart | The fields are non-nullable and their constructor parameters are `required`. | Pass them when you build a model; `!` and `?.` on them become analyzer warnings. |
| Ruby | No signature changes. | `Model.new` without a required attribute raises `ArgumentError` (`id cannot be nil`). |

If an answer ever lacked a required field, the Python, Go, .NET, Dart and Ruby SDKs would refuse the whole answer with an
error (Dart, outside debug mode, reads a missing list as empty), and TypeScript would for a missing list; Java would read it
and leave the field null, and PHP would read it but throw a `TypeError` when that field is read.

**Go: a field the SDK does not know no longer breaks an answer.** 98 Go models refused an answer that carried a property
they did not know, so a field added to the API could break a deployed client. Every Go model now ignores unknown
properties, as the other SDKs do.

Compared with 0.1.0-rc.6, the API description changed in 75 breaking and 80 other ways, most of them the move from
`/api/v1/b2b/...` to `/v1/...`. The sections above say what each change means for your code.

## 0.1.0-rc.6

**Breaking: type names.** The API now publishes its types under the names partners read, and every SDK uses them as they
are, so a type has the same name in the API reference and in each SDK. The names are shorter. Nothing on the wire changes:
the routes, the JSON and the client methods are the same. Code that names a type (an import, a type annotation,
`new HotelRoomGuestsInput(...)`) needs the new name. So does code that reads a union by its variant's name: in Java
`getOfferFixedPrice()` is now `getFixedPrice()` (C#: `GetFixedPrice()`), and in Go the field `OfferFixedPrice` is
`FixedPrice` and `EsimOrderFulfillmentUpdatedWebhookItem` is `EsimFulfillmentUpdate`. Dart's `fixed` and `range` stay.

- A product prefix stays only where the type is that product's: `FlightAirline` is now `Airline` and `HotelRoom` is `Room`,
  while `HotelOffer` keeps its prefix.
- Words that repeated the type around it are gone: `HotelPurchasedProductGuest` is `Guest`, `EsimProductOptionApn` is
  `EsimApn`.
- Money and offers: `MoneyAmount` is `Money` and `MoneyAmountInput` is `MoneyInput`. An offer's price is a `FixedPrice` or a
  `PriceRange`, and its value a `FixedValue` or a `ValueRange`.
- The order of several products: `OrderRequest` is `MixedOrderRequest` and `OrderQuote` is `MixedQuote`.
- The fulfillment webhook: `OrderFulfillmentUpdatedWebhook` is `OrderFulfillmentWebhook`, and its items are
  `FulfillmentUpdate` (`HotelFulfillmentUpdate` and the other products).

Every renamed type:

| Before | Now |
|---|---|
| `AirtimeFulfillmentRedemption` | `AirtimeRedemption` |
| `AirtimeFulfillmentResultDetails` | `AirtimeFulfillmentDetails` |
| `AirtimeNumberCarrier` | `NumberCarrier` |
| `AirtimeNumberLookup` | `NumberLookup` |
| `AirtimeOrderFulfillmentUpdatedWebhookItem` | `AirtimeFulfillmentUpdate` |
| `AirtimePhoneLookup` | `CarrierDetection` |
| `AirtimePurchasedProductDetails` | `AirtimePurchasedDetails` |
| `AirtimePurchasedProductInputs` | `AirtimePurchasedInputs` |
| `AmendmentListItem` | `Amendment` |
| `BookingTermsCancellation` | `CancellationTerms` |
| `BookingTermsCancellationStep` | `CancellationStep` |
| `BookingTermsChangeFee` | `ChangeFee` |
| `BookingTermsChanges` | `ChangeTerms` |
| `BookingTermsPayment` | `PaymentTerms` |
| `CancellationPenaltyDetail` | `CancellationPenalty` |
| `CompatibleDevice` | `EsimCompatibleDevice` |
| `EsimFulfillmentInstallation` | `EsimInstallation` |
| `EsimFulfillmentResultDetails` | `EsimFulfillmentDetails` |
| `EsimOfferData` | `EsimDataAllowance` |
| `EsimOrderFulfillmentUpdatedWebhookItem` | `EsimFulfillmentUpdate` |
| `EsimProductOptionApn` | `EsimApn` |
| `EsimProductOptionData` | `EsimProfileData` |
| `EsimProductOptionProperties` | `EsimFeatures` |
| `EsimProductOptionRestriction` | `EsimRestrictions` |
| `EsimProductOptionSms` | `EsimSmsAllowance` |
| `EsimProductOptionVoice` | `EsimVoiceAllowance` |
| `EsimProvisionRequest` | `EsimProvisioningRequest` |
| `EsimPurchasedProductDestination` | `EsimPurchasedDestination` |
| `EsimPurchasedProductDetails` | `EsimPurchasedDetails` |
| `EsimPurchasedProductInputs` | `EsimPurchasedInputs` |
| `EsimSendSmsRequest` | `EsimSmsRequest` |
| `EsimSendSmsResult` | `EsimSmsResult` |
| `EsimSim` | `SimCard` |
| `FlightAircraft` | `Aircraft` |
| `FlightAirline` | `Airline` |
| `FlightAirport` | `Airport` |
| `FlightAvailableService` | `AvailableService` |
| `FlightBaggage` | `Baggage` |
| `FlightBookedAirline` | `BookedAirline` |
| `FlightBookedAirport` | `BookedAirport` |
| `FlightBookedBaggage` | `BookedBaggage` |
| `FlightBookedSegment` | `BookedSegment` |
| `FlightBookedSlice` | `BookedSlice` |
| `FlightBookedStop` | `BookedStop` |
| `FlightFare` | `Fare` |
| `FlightFareOption` | `FareOption` |
| `FlightFareSegment` | `SegmentFare` |
| `FlightFulfillmentBooking` | `FlightBooking` |
| `FlightFulfillmentResultDetails` | `FlightFulfillmentDetails` |
| `FlightIdentityDocument` | `IdentityDocument` |
| `FlightIdentityDocumentInput` | `IdentityDocumentInput` |
| `FlightItinerary` | `Itinerary` |
| `FlightItineraryDetails` | `ItineraryDetails` |
| `FlightOfferPassenger` | `OfferPassenger` |
| `FlightOrderFulfillmentUpdatedWebhookItem` | `FlightFulfillmentUpdate` |
| `FlightPassenger` | `Passenger` |
| `FlightPassengerInput` | `PassengerInput` |
| `FlightPassengerPricing` | `PassengerPricing` |
| `FlightPriceBreakdown` | `PriceBreakdown` |
| `FlightPriceLine` | `PriceLine` |
| `FlightPurchasedProductDetails` | `FlightPurchasedDetails` |
| `FlightPurchasedProductInputs` | `FlightPurchasedInputs` |
| `FlightSeatMap` | `SeatMap` |
| `FlightSeatMapAvailability` | `SeatAvailability` |
| `FlightSeatMapElement` | `SeatMapElement` |
| `FlightSeatMapRow` | `SeatMapRow` |
| `FlightSeatMapSection` | `SeatMapSection` |
| `FlightSeatMapSegment` | `SeatMapSegment` |
| `FlightSegment` | `Segment` |
| `FlightSegmentReference` | `SegmentPosition` |
| `FlightSelectedSeat` | `SelectedSeat` |
| `FlightSelectedSeatInput` | `SelectedSeatInput` |
| `FlightSelectedService` | `SelectedService` |
| `FlightSelectedServiceInput` | `SelectedServiceInput` |
| `FlightServiceBaggage` | `BaggageLimits` |
| `FlightSlice` | `Slice` |
| `FlightStop` | `Stop` |
| `FlightTicketDocument` | `TicketDocument` |
| `GeoCoordinates` | `Coordinates` |
| `GiftCardFulfillmentDelivery` | `GiftCardDelivery` |
| `GiftCardFulfillmentResultDetails` | `GiftCardFulfillmentDetails` |
| `GiftCardIssuedCode` | `GiftCardCode` |
| `GiftCardOrderFulfillmentUpdatedWebhookItem` | `GiftCardFulfillmentUpdate` |
| `GiftCardPurchasedProductDetails` | `GiftCardPurchasedDetails` |
| `GiftCardPurchasedProductInputs` | `GiftCardPurchasedInputs` |
| `HotelAddress` | `Address` |
| `HotelAmenity` | `Amenity` |
| `HotelAmenityCatalog` | `AmenityCatalog` |
| `HotelAmenityFacet` | `AmenityFacet` |
| `HotelAmenityGroupFacet` | `AmenityGroup` |
| `HotelBed` | `Bed` |
| `HotelBoardType` | `MealPlan` |
| `HotelBookingGuestInput` | `GuestInput` |
| `HotelCatalogProperty` | `CatalogHotel` |
| `HotelDestinationSummary` | `HotelDestination` |
| `HotelDisplayGroup` | `UnitAvailability` |
| `HotelFulfillmentBooking` | `HotelBooking` |
| `HotelFulfillmentResultDetails` | `HotelFulfillmentDetails` |
| `HotelGuestCapabilities` | `GuestAmendmentCapabilities` |
| `HotelGuestNameRules` | `GuestNameRules` |
| `HotelHolderInput` | `LeadGuestInput` |
| `HotelLocationSearch` | `HotelSearchArea` |
| `HotelNightlyPrice` | `NightlyPrice` |
| `HotelOrderFulfillmentUpdatedWebhookItem` | `HotelFulfillmentUpdate` |
| `HotelProperty` | `Hotel` |
| `HotelPropertyDetails` | `HotelDetails` |
| `HotelPropertySearchDetails` | `HotelSearchDetails` |
| `HotelPropertySearchResult` | `HotelSearchResult` |
| `HotelPurchasedProductDetails` | `HotelPurchasedDetails` |
| `HotelPurchasedProductGuest` | `Guest` |
| `HotelPurchasedProductHolder` | `LeadGuest` |
| `HotelPurchasedProductHotel` | `BookedHotel` |
| `HotelPurchasedProductInputs` | `HotelPurchasedInputs` |
| `HotelPurchasedProductRoomGuests` | `RoomGuests` |
| `HotelQuotedRate` | `QuotedRate` |
| `HotelRoom` | `Room` |
| `HotelRoomCapabilities` | `RoomAmendmentCapabilities` |
| `HotelRoomGuestsInput` | `RoomGuestsInput` |
| `HotelRoomOccupancy` | `RoomOccupancy` |
| `HotelSearchCoordinates` | `CoordinatesInput` |
| `HotelSearchRoom` | `RoomOccupancyInput` |
| `HotelSpecialRequestOption` | `SpecialRequestOption` |
| `HotelSurcharge` | `Surcharge` |
| `HotelTaxBreakdown` | `TaxBreakdown` |
| `HotelTaxItem` | `Tax` |
| `MoneyAmount` | `Money` |
| `MoneyAmountInput` | `MoneyInput` |
| `OfferFixedPrice` | `FixedPrice` |
| `OfferFixedValue` | `FixedValue` |
| `OfferRangePrice` | `PriceRange` |
| `OfferRangeValue` | `ValueRange` |
| `OrderFulfillmentUpdatedWebhook` | `OrderFulfillmentWebhook` |
| `OrderFulfillmentUpdatedWebhookData` | `OrderFulfillmentWebhookData` |
| `OrderFulfillmentUpdatedWebhookItem` | `FulfillmentUpdate` |
| `OrderItemManualAmendmentStatusWebhook` | `ManualAmendmentStatusWebhook` |
| `OrderItemManualAmendmentStatusWebhookData` | `ManualAmendmentStatusWebhookData` |
| `OrderQuote` | `MixedQuote` |
| `OrderQuoteItem` | `MixedQuoteItem` |
| `OrderRequest` | `MixedOrderRequest` |
| `OrderRequestItem` | `MixedOrderRequestItem` |
| `RequiredInput` | `InputField` |
| `UtilityBiller` | `Biller` |
| `UtilityBillerDetails` | `BillerDetails` |
| `UtilityFulfillmentResultDetails` | `UtilityFulfillmentDetails` |
| `UtilityOrderFulfillmentUpdatedWebhookItem` | `UtilityFulfillmentUpdate` |
| `UtilityPurchasedProductDetails` | `UtilityPurchasedDetails` |
| `UtilityPurchasedProductInputs` | `UtilityPurchasedInputs` |

**Breaking** (108)
- `GET /api/v1/b2b/airtime/operators/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/airtime/operators/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/airtime/operators/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/airtime/operators/{id}`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `offers/items/value` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/airtime/operators/{id}`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `offers/items/value` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/airtime/operators/{id}`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `offers/items/value` response property `oneOf` list for the response status `200` (media type: text/plain)
- `POST /api/v1/b2b/airtime/order`: added `#/components/schemas/AirtimeFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: application/json)
- `POST /api/v1/b2b/airtime/order`: added `#/components/schemas/AirtimeFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: text/json)
- `POST /api/v1/b2b/airtime/order`: added `#/components/schemas/AirtimeFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/airtime/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/airtime/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/airtime/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/airtime/search`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `items/items/offers/items/value` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/airtime/search`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `items/items/offers/items/value` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/airtime/search`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `items/items/offers/items/value` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/esims/catalog`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/esims/catalog`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/esims/catalog`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/esims/destinations/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/esims/destinations/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/esims/destinations/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `POST /api/v1/b2b/esims/order`: added `#/components/schemas/EsimFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: application/json)
- `POST /api/v1/b2b/esims/order`: added `#/components/schemas/EsimFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: text/json)
- `POST /api/v1/b2b/esims/order`: added `#/components/schemas/EsimFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/esims/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/esims/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/esims/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `POST /api/v1/b2b/flights/order`: added `#/components/schemas/FlightFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: application/json)
- `POST /api/v1/b2b/flights/order`: added `#/components/schemas/FlightFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: text/json)
- `POST /api/v1/b2b/flights/order`: added `#/components/schemas/FlightFulfillmentResult` to the `items/items/fulfillment/result` response property `oneOf` list for the response status `200` (media type: text/plain)
- `POST /api/v1/b2b/flights/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `POST /api/v1/b2b/flights/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `POST /api/v1/b2b/flights/search`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/giftcards/cards/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/giftcards/cards/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/giftcards/cards/{id}`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `offers/items/price` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/giftcards/cards/{id}`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `offers/items/value` response property `oneOf` list for the response status `200` (media type: application/json)
- `GET /api/v1/b2b/giftcards/cards/{id}`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `offers/items/value` response property `oneOf` list for the response status `200` (media type: text/json)
- `GET /api/v1/b2b/giftcards/cards/{id}`: added `#/components/schemas/FixedValue, #/components/schemas/ValueRange` to the `offers/items/value` response property `oneOf` list for the response status `200` (media type: text/plain)
- `GET /api/v1/b2b/giftcards/catalog`: added `#/components/schemas/FixedPrice, #/components/schemas/PriceRange` to the `items/items/offers/items/price` response property `oneOf` list for the response status `200` (media type: application/json)
- ...and 68 more

**Other API changes** (649)
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Currencies.CurrencyDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.B2BAmendmentPaymentDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.B2BAmendmentQuoteOptionDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.B2BAmendmentRefundDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.B2BOrderItemAmendmentDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.B2BPartnerAmendmentItemDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.B2BPartnerAmendmentListDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.ConfirmB2BAmendmentRequestDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.ConfirmB2BAmendmentResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Amendments.B2B.QuoteB2BAmendmentResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Refunds.B2B.B2BCancellationPreviewResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Refunds.B2B.B2BOrderItemRefundResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Refunds.CancellationPenaltyDetailDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.OrderItems.Updates.B2B.UpdateB2BOrderItemRequestDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.B2BOrderAcceptedResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.B2BOrderFulfillmentStatusResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.B2BOrderListResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.B2BOrderResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.B2BOrderSummaryResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.RetryB2BOrderFulfillmentResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.B2B.RetryB2BOrderPaymentResponseDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Fulfillment.PartnerFulfillmentFailureDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Fulfillment.PartnerFulfillmentProgressDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.IPartnerFulfillmentResultDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.IPartnerPurchasedProductDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerActionDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderActionsDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderCancellationPreviewDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderCancellationPreviewTotalsDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemActionsDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemCancellationPreviewDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto.Airtime`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto.Esim`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto.Flight`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto.GiftCard`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto.Hotel`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemDto.Utility`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderItemRefundDto`
- removed the schema `OpenRoamBackend.Application.Shared.Commerce.DTOs.Orders.Partner.PartnerOrderOriginDto`
- ...and 609 more

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
