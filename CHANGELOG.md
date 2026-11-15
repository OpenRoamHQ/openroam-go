# Changelog

All OpenRoam SDKs share one version and one changelog. Entries are written by the release process from the
changes in the Partner API description.

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
