package openroam

import (
	"encoding/json"
	"testing"
)

// A product's own fields sit on the item and on the quote item, beside the common keys. Offers and fare options keep
// `details`. The payloads are the examples of the API description, shortened to one entry per list.
var flatPayloads = map[string]string{
	"airtimeSearch": `
{"items":[{"id":"operator_tq7vbm2xk4hdn6wzr3pyjc5fae","name":"Bakcell Azerbaijan","countryCode":"AZ","category":null,
"logoUrl":"https://assets.openroam.io/services/image.png","description":null,
"offers":[{"publicId":"airtime_raanuxyz6cxzie3w3my6c3r2wm","name":"Bakcell Azerbaijan 40 AZN",
"description":"Top-up for **Bakcell Azerbaijan** numbers.","value":{"type":"fixed","amount":40,"currency":"AZN"},
"price":{"amount":25.0,"currency":"USD"},"taxes":{"treatment":"notStated","amount":null,"lines":[]},
"requiredInputs":[{"name":"phoneNumber","type":"phone","required":true,
"description":"The recipient's mobile number in international format, starting with + or 00.",
"labelKey":"inputs.phoneNumber"}],"details":{"redemptionType":"directTopUp",
"readMore":"Credit arrives within minutes."}}]}],"pagination":{"currentCount":1,"totalCount":1,"skip":0,"limit":10,
"hasMore":false,"page":1,"totalPages":1},"phoneLookup":{"status":"detected","carrierName":"Bakcell LLC",
"countryCode":"AZ"}}`,
	"airtimeQuote": `
{"currency":"USD","items":[{"publicId":"airtime_raanuxyz6cxzie3w3my6c3r2wm","quantity":1,
"name":"Bakcell Azerbaijan 40 AZN","value":{"amount":40,"currency":"AZN"},"price":{"amount":25.0,"currency":"USD"},
"taxes":{"treatment":"notStated","amount":null,"lines":[]},"estimatedDeliveredValue":{"amount":40,"currency":"AZN"}}],
"total":{"amount":25.0,"currency":"USD"}}`,
	"esimSearch": `
{"items":[{"id":"destn_k3m7q2t5w4y6b2d4f6h3j5n7pq","name":"Europe","countryCode":null,"category":"regional",
"logoUrl":"https://cdn.openroam.test/esims/europe.png","description":null,
"offers":[{"publicId":"esim_k7m2q4t6v3x5z2b4d6f3h5jn2p","name":"Europe 3 GB 7 days","description":null,
"price":{"amount":25.0,"currency":"USD"},"taxes":{"treatment":"notStated","amount":null,"lines":[]},
"requiredInputs":[{"name":"isInstant","type":"boolean","required":false,
"description":"Defaults to true, which creates the eSIM when the order is fulfilled. False delivers an entitlement, and the eSIM is created later through provisioning.",
"labelKey":"inputs.isInstant"}],"details":{"validityDays":7,"data":{"amount":3072,"unit":"MB","isUnlimited":false,
"label":"3 GB"},"voice":null,"sms":null,"apn":{"apn":"internet","apnUsername":null,"apnPassword":null},
"properties":{"withKyc":false,"withTopUp":true,"withHotspot":true,"withDataRoaming":true,"withSubscription":false,
"withDestinationInstall":false,"withUsageCheck":true,"technology":"eSIM","activation":"first-use",
"charge":"one-time"},"restrictions":{"throttle":false,"throttleSpeed":null,"throttleUnit":null,
"throttleThreshold":null,"throttleThresholdUnit":null,"expireInDays":30},
"notes":["Data only. Calls and SMS are not included."],"plan":null}}],"destinationCode":"EU","region":"Europe",
"coverage":[{"countryCode":"DE","name":"Germany","networks":["4G"],"carriers":["Telekom"]}]}],
"pagination":{"currentCount":1,"totalCount":1,"skip":0,"limit":10,"hasMore":false,"page":1,"totalPages":1}}`,
	"giftCardSearch": `
{"items":[{"id":"card_tzkc6b2xcqznidhkouwh2rz7wi","name":"Amazon","countryCode":"US","category":"shopping",
"logoUrl":"https://assets.openroam.io/services/image.png","description":"Redeem at amazon.com",
"offers":[{"publicId":"giftcard_p5udmgvqke2bcrwlk35drsqfqm","name":"Amazon 50 USD","description":null,
"value":{"type":"fixed","amount":50,"currency":"USD"},"price":{"amount":48.98,"currency":"USD"},
"taxes":{"treatment":"notStated","amount":null,"lines":[]},"requiredInputs":[{"name":"recipientEmail","type":"email",
"required":false,"description":"Sends the card to this email address as a gift.","labelKey":"inputs.recipientEmail",
"maxLength":254}]}],"isRedeemableGlobally":false,
"redemptionInstructions":"Go to amazon.com/redeem and enter the code."}],"pagination":{"currentCount":1,
"totalCount":1,"skip":0,"limit":10,"hasMore":false,"page":1,"totalPages":1}}`,
	"utilitySearch": `
{"items":[{"id":"biller_3jx464dkw4ekdjgj4qxuib5wki","name":"Azeriqaz","countryCode":"AZ","category":"gas",
"logoUrl":null,"description":null,"offers":[{"publicId":"utility_7tq2mw5kxr3nbv6hzc4dpe2ysa","name":"Azeriqaz 40 AZN",
"description":null,"value":{"type":"fixed","amount":40,"currency":"AZN"},"price":{"amount":27.5,"currency":"USD"},
"taxes":{"treatment":"notStated","amount":null,"lines":[]},"requiredInputs":[{"name":"accountNumber","type":"text",
"required":true,"description":"The customer's account or meter number with the biller.",
"labelKey":"inputs.accountNumber","maxLength":128}]}],"serviceType":"postpaid","deliveryType":"directCredit"}],
"pagination":{"currentCount":1,"totalCount":1,"skip":0,"limit":10,"hasMore":false,"page":1,"totalPages":1}}`,
	"utilityQuote": `
{"currency":"USD","items":[{"publicId":"utility_7tq2mw5kxr3nbv6hzc4dpe2ysa","quantity":1,"name":"Azeriqaz 40 AZN",
"value":{"amount":40,"currency":"AZN"},"price":{"amount":27.5,"currency":"USD"},"taxes":{"treatment":"notStated",
"amount":null,"lines":[]},"deliveryType":"directCredit"}],"total":{"amount":27.5,"currency":"USD"}}`,
	"hotelSearch": `
{"items":[{"id":"hotel_q7m2k4t6v3x5z2b4d6f3h5jn2p","name":"Hotel Arts Barcelona","countryCode":"ES","category":null,
"logoUrl":"https://assets.openroam.io/services/image.png",
"description":"<p>A seafront tower in the Olympic Port.</p>","offers":[{"publicId":"rate_x5z2b4d6f3h5jn2pq7m2k4t6v3",
"name":"Deluxe King, Breakfast, Free cancellation","description":null,"price":{"amount":642.1,"currency":"USD"},
"taxes":{"treatment":"included","amount":{"amount":41.2,"currency":"USD"},"lines":[{"code":null,"codeSystem":null,
"name":"Tax or fee","kind":"unknown","payableAt":"prepaid","per":null,"jurisdiction":null,"amount":{"amount":38.0,
"currency":"EUR"},"approxAmount":{"amount":41.2,"currency":"USD"}}]},"details":{"expiresAt":"2026-12-01T10:30:00Z",
"rooms":[{"roomGroupKey":"4c8e3ae1d379f759bc372243","name":"Deluxe King","description":null,"beds":[],
"maxOccupancy":3,"sizeSquareMeters":null,"amenities":[],"imageUrls":[],"roomIndex":1,"adultsCount":2,
"childrenCount":1,"childrenAges":[7]}],"board":{"code":"BB","name":"Bed and breakfast","boardType":"breakfast"},
"terms":{"payment":{"type":"payNow","dueNow":{"amount":642.1,"currency":"USD"},"dueLater":null,"dueLaterAt":null},
"cancellation":{"type":"fullyRefundable","freeUntil":"2026-12-08T23:59:00Z",
"freeUntilLocalTime":"2026-12-09T00:59:00","timezone":"Europe/Madrid","schedule":[{"from":"2026-12-08T23:59:00Z",
"fromLocalTime":"2026-12-09T00:59:00","refund":{"amount":428.07,"currency":"USD"},"penalty":{"amount":214.03,
"currency":"USD"}}],"fee":null},"changes":{"allowed":true,"fee":null}},"priceIsBinding":true,
"nightlyPrices":[{"date":"2026-12-10","amount":214.03,"currency":"USD"}],"surcharges":[],"rateTags":[],"remarks":[],
"withTransfers":false,"rateRules":[],"amendabilityStatus":"yes","guestNameRules":{"minLength":1,"maxLength":50,
"pattern":"^(?=.*[a-zA-Z])[a-zA-Z\\s\\-'\\.]+$",
"allowedCharacters":"letters, spaces, hyphens, apostrophes and periods","distinctFullNames":false,
"guestNamesOptional":false},"minimumSellingPrice":null}}],"starRating":5,
"address":{"line1":"Carrer de la Marina 19-21","line2":null,"city":"Barcelona","stateProvince":null,
"postalCode":"08005","countryCode":"ES","countryName":null},"coordinates":{"latitude":41.386,"longitude":2.196},
"imageUrls":["https://assets.openroam.io/services/image.png"],"amenityGroups":[{"key":"PropertyAmenities",
"codes":["FreeWifi"]}],"contact":null,"displayGroup":null}],"pagination":{"currentCount":1,"limit":1,"hasMore":false},
"nextCursor":null,"availableFilters":{"amenities":[{"code":"FreeWifi","label":"Free Wi-Fi","count":1}],
"amenityGroups":[{"key":"PropertyAmenities","label":"Property amenities","codes":["FreeWifi"]}]}}`,
	"catalogHotel": `
{"id":"hotel_q7m2k4t6v3x5z2b4d6f3h5jn2p","name":"Hotel Arts Barcelona","countryCode":"ES","category":null,
"logoUrl":"https://assets.openroam.io/services/image.png",
"description":"<p>A seafront tower in the Olympic Port.</p>","starRating":5,
"address":{"line1":"Carrer de la Marina 19-21","line2":null,"city":"Barcelona","stateProvince":null,
"postalCode":"08005","countryCode":"ES","countryName":null},"coordinates":{"latitude":41.386,"longitude":2.196},
"imageUrls":["https://assets.openroam.io/services/image.png"],"amenityGroups":[{"key":"PropertyAmenities",
"codes":["FreeWifi"]}],"contact":null,"displayGroup":null}`,
	"hotelQuote": `
{"currency":"USD","items":[{"publicId":"rate_b4d6f3h5jn2pq7m2k4t6v3x5z2","quantity":1,
"name":"Hotel Arts Barcelona \u00b7 Deluxe King, Breakfast, Free cancellation","price":{"amount":648.9,
"currency":"USD"},"taxes":{"treatment":"included","amount":{"amount":41.64,"currency":"USD"},"lines":[{"code":null,
"codeSystem":null,"name":"Tax or fee","kind":"unknown","payableAt":"prepaid","per":null,"jurisdiction":null,
"amount":{"amount":38.0,"currency":"EUR"},"approxAmount":{"amount":41.64,"currency":"USD"}}]},
"expiresAt":"2026-12-01T10:42:00Z","availability":"confirmed","checkInDate":"2026-12-10","checkOutDate":"2026-12-13",
"requestedRooms":[{"adults":2,"children":1,"childrenAges":[7]}],"terms":{"payment":{"type":"payNow",
"dueNow":{"amount":648.9,"currency":"USD"},"dueLater":null,"dueLaterAt":null},
"cancellation":{"type":"fullyRefundable","freeUntil":"2026-12-08T23:59:00Z",
"freeUntilLocalTime":"2026-12-09T00:59:00","timezone":"Europe/Madrid","schedule":[{"from":"2026-12-08T23:59:00Z",
"fromLocalTime":"2026-12-09T00:59:00","refund":{"amount":432.6,"currency":"USD"},"penalty":{"amount":216.3,
"currency":"USD"}}],"fee":null},"changes":{"allowed":true,"fee":null}},
"rate":{"rooms":[{"roomGroupKey":"4c8e3ae1d379f759bc372243","name":"Deluxe King","description":null,"beds":[],
"maxOccupancy":3,"sizeSquareMeters":null,"amenities":[],"imageUrls":[],"roomIndex":1,"adultsCount":2,
"childrenCount":1,"childrenAges":[7]}],"board":{"code":"BB","name":"Bed and breakfast","boardType":"breakfast"},
"priceIsBinding":true,"nightlyPrices":[{"date":"2026-12-10","amount":216.3,"currency":"USD"}],"rateTags":[],
"remarks":[],"rateConditions":[],"withTransfers":false,"rateRules":[],"amendabilityStatus":"yes","surcharges":[],
"minimumSellingPrice":null}}],"total":{"amount":648.9,"currency":"USD"}}`,
	"flightSearch": `
{"items":[{"id":"itin_dd6d0153af8985bfc97cc5707e58af90",
"name":"LHR \u21c4 JFK \u00b7 12\u201319 Dec \u00b7 British Airways","countryCode":"US","category":null,
"logoUrl":"https://assets.openroam.io/services/image.png","description":null,
"offers":[{"publicId":"flight_dwqbtkmavshyh6pll7ixbouz2o","name":"Economy Basic, Non-refundable","description":null,
"price":{"amount":1184.2,"currency":"USD"},"taxes":{"treatment":"included","amount":{"amount":204.2,"currency":"USD"},
"lines":[]},"details":{"expiresAt":"2026-12-01T10:30:00Z","cabinClass":"economy","fare":{"fareBasisCode":"OXZ0RO",
"brandName":"Basic","fareFamily":"Economy Basic","validatingCarrier":"BA"},"terms":{"payment":{"type":"payNow",
"dueNow":{"amount":1184.2,"currency":"USD"},"dueLater":null,"dueLaterAt":null},"cancellation":{"type":"nonRefundable",
"freeUntil":null,"freeUntilLocalTime":null,"timezone":"Europe/London","schedule":[],"fee":null},
"changes":{"allowed":true,"fee":null}},"includedBaggage":[{"type":"carryOn","quantity":1,"weightKg":null,
"description":null}],"passengerPricing":[{"passengerType":"adult","count":1,"pricePerPassenger":null,
"fareBasisCode":"OXZ0RO","bookingClass":"O"}],"offerPassengers":[{"passengerIndex":1,"type":"adult","age":null,
"fareType":null,"includedBaggage":[{"type":"carryOn","quantity":1,"weightKg":null,"description":null}],
"availableServices":[]}],"availableServices":[],"identityDocumentsRequired":true,
"supportedDocumentTypes":["passport"],"segments":[{"sliceIndex":1,"segmentIndex":1,"cabinClass":"economy",
"baggage":[{"type":"carryOn","quantity":1,"weightKg":null,"description":null}]}]}}],"owner":{"iataCode":"BA",
"name":"British Airways","logoUrl":"https://assets.openroam.io/services/image.png"},
"slices":[{"origin":{"iataCode":"LHR","icaoCode":"EGLL","name":"Heathrow Airport","city":"London",
"cityIataCode":"LON","countryCode":"GB","terminal":"5","timezone":"Europe/London"},"destination":{"iataCode":"JFK",
"icaoCode":"KJFK","name":"John F. Kennedy International Airport","city":"New York","cityIataCode":"NYC",
"countryCode":"US","terminal":"8","timezone":"America/New_York"},"departureAt":"2026-12-12T08:30:00Z",
"departureLocalTime":"2026-12-12T08:30:00","arrivalAt":"2026-12-12T16:20:00Z",
"arrivalLocalTime":"2026-12-12T11:20:00","durationMinutes":470,"stops":0,"segments":[{"origin":{"iataCode":"LHR",
"icaoCode":"EGLL","name":"Heathrow Airport","city":"London","cityIataCode":"LON","countryCode":"GB","terminal":"5",
"timezone":"Europe/London"},"destination":{"iataCode":"JFK","icaoCode":"KJFK",
"name":"John F. Kennedy International Airport","city":"New York","cityIataCode":"NYC","countryCode":"US",
"terminal":"8","timezone":"America/New_York"},"departureAt":"2026-12-12T08:30:00Z",
"departureLocalTime":"2026-12-12T08:30:00","arrivalAt":"2026-12-12T16:20:00Z",
"arrivalLocalTime":"2026-12-12T11:20:00","durationMinutes":470,"marketingCarrier":{"iataCode":"BA",
"name":"British Airways","logoUrl":"https://assets.openroam.io/services/image.png"},"marketingFlightNumber":"117",
"operatingCarrier":{"iataCode":"BA","name":"British Airways",
"logoUrl":"https://assets.openroam.io/services/image.png"},"operatingFlightNumber":"117","aircraft":{"iataCode":"77W",
"name":"Boeing 777-300ER","imageUrl":null},"stops":[],"isCodeshare":false}]}],"marketingCarriers":[{"iataCode":"BA",
"name":"British Airways","logoUrl":"https://assets.openroam.io/services/image.png"}],
"operatingCarriers":[{"iataCode":"BA","name":"British Airways",
"logoUrl":"https://assets.openroam.io/services/image.png"}],"totalDurationMinutes":885,"totalStops":0,
"totalEmissionsKg":1236}],"pagination":{"currentCount":1,"totalCount":1,"skip":0,"limit":10,"hasMore":false,"page":1,
"totalPages":1}}`,
	"flightPlaces": `
{"items":[{"id":"LON","name":"London","countryCode":"GB","category":"city","logoUrl":null,"iataCode":"LON",
"cityName":"London","cityIataCode":"LON","timezone":"Europe/London","airports":[{"iataCode":"LGW",
"name":"Gatwick Airport","timezone":"Europe/London"}]}]}`,
	"flightQuote": `
{"currency":"USD","items":[{"publicId":"flight_guis4oiryxjtxfalagepbfdcpk","quantity":1,
"name":"LHR \u21c4 JFK \u00b7 12\u201319 Dec \u00b7 Economy Basic, Non-refundable","price":{"amount":1184.2,
"currency":"USD"},"taxes":{"treatment":"included","amount":{"amount":204.2,"currency":"USD"},"lines":[{"code":"US",
"codeSystem":"iata","name":"Transportation tax","kind":"tax","payableAt":"prepaid","per":null,"jurisdiction":"US",
"amount":{"amount":47.0,"currency":"USD"},"approxAmount":null}]},"expiresAt":"2026-12-01T10:30:00Z",
"quotedAt":"2026-12-01T10:05:12Z","hasSeatMaps":true,"terms":{"payment":{"type":"payNow","dueNow":{"amount":1184.2,
"currency":"USD"},"dueLater":null,"dueLaterAt":null},"cancellation":{"type":"nonRefundable","freeUntil":null,
"freeUntilLocalTime":null,"timezone":"Europe/London","schedule":[],"fee":null},"changes":{"allowed":true,"fee":null}},
"includedBaggage":[{"type":"carryOn","quantity":1,"weightKg":null,"description":null}],
"availableServices":[{"id":"flight_4224ss6vukylqlb2iphmrx3u52","type":"baggage","price":{"amount":30.0,
"currency":"USD"},"passengerIndexes":[1],"segments":[{"sliceIndex":1,"segmentIndex":1}],"maxQuantity":2,
"baggage":{"type":"checked","maximumWeightKg":23,"maximumLengthCm":null,"maximumHeightCm":null,
"maximumDepthCm":null}}],"fareOptions":[{"publicId":"flight_fi6uxnnuusc37omsl7vrs6z6lt",
"name":"Economy Basic, Non-refundable","price":{"amount":1184.2,"currency":"USD"},"taxes":{"treatment":"included",
"amount":{"amount":204.2,"currency":"USD"},"lines":[{"code":"US","codeSystem":"iata","name":"Transportation tax",
"kind":"tax","payableAt":"prepaid","per":null,"jurisdiction":"US","amount":{"amount":47.0,"currency":"USD"},
"approxAmount":null}]},"details":{"expiresAt":"2026-12-01T10:30:00Z","cabinClass":"economy",
"fare":{"fareBasisCode":"OXZ0RO","brandName":"Basic","fareFamily":"Economy Basic","validatingCarrier":"BA"},
"terms":{"payment":{"type":"payNow","dueNow":{"amount":1184.2,"currency":"USD"},"dueLater":null,"dueLaterAt":null},
"cancellation":{"type":"nonRefundable","freeUntil":null,"freeUntilLocalTime":null,"timezone":"Europe/London",
"schedule":[],"fee":null},"changes":{"allowed":true,"fee":null}},"includedBaggage":[{"type":"carryOn","quantity":1,
"weightKg":null,"description":null}],"passengerPricing":[{"passengerType":"adult","count":1,"pricePerPassenger":null,
"fareBasisCode":"OXZ0RO","bookingClass":"O"}],"offerPassengers":[{"passengerIndex":1,"type":"adult","age":null,
"fareType":null,"includedBaggage":[{"type":"carryOn","quantity":1,"weightKg":null,"description":null}],
"availableServices":[{"id":"flight_4224ss6vukylqlb2iphmrx3u52","type":"baggage","price":{"amount":30.0,
"currency":"USD"},"passengerIndexes":[1],"segments":[{"sliceIndex":1,"segmentIndex":1}],"maxQuantity":2,
"baggage":{"type":"checked","maximumWeightKg":23,"maximumLengthCm":null,"maximumHeightCm":null,
"maximumDepthCm":null}}]}],"availableServices":[{"id":"flight_4224ss6vukylqlb2iphmrx3u52","type":"baggage",
"price":{"amount":30.0,"currency":"USD"},"passengerIndexes":[1],"segments":[{"sliceIndex":1,"segmentIndex":1}],
"maxQuantity":2,"baggage":{"type":"checked","maximumWeightKg":23,"maximumLengthCm":null,"maximumHeightCm":null,
"maximumDepthCm":null}}],"identityDocumentsRequired":true,"supportedDocumentTypes":["passport"],
"segments":[{"sliceIndex":1,"segmentIndex":1,"cabinClass":"economy","baggage":[{"type":"carryOn","quantity":1,
"weightKg":null,"description":null}]}]}}]}],"total":{"amount":1184.2,"currency":"USD"}}`,
	"mixedQuote": `
{"currency":"USD","items":[{"productType":"airtime","publicId":"airtime_raanuxyz6cxzie3w3my6c3r2wm","quantity":1,
"name":"Bakcell Azerbaijan 40 AZN","value":{"amount":40,"currency":"AZN"},"price":{"amount":25.0,"currency":"USD"},
"taxes":{"treatment":"notStated","amount":null,"lines":[]},"estimatedDeliveredValue":{"amount":40,"currency":"AZN"}},
{"productType":"utility","publicId":"utility_7tq2mw5kxr3nbv6hzc4dpe2ysa","quantity":1,"name":"Azeriqaz 40 AZN",
"value":{"amount":40,"currency":"AZN"},"price":{"amount":27.5,"currency":"USD"},"taxes":{"treatment":"notStated",
"amount":null,"lines":[]},"deliveryType":"directCredit"},{"productType":"hotel",
"publicId":"rate_b4d6f3h5jn2pq7m2k4t6v3x5z2","quantity":1,
"name":"Hotel Arts Barcelona \u00b7 Deluxe King, Breakfast, Free cancellation","price":{"amount":648.9,
"currency":"USD"},"taxes":{"treatment":"included","amount":{"amount":41.64,"currency":"USD"},"lines":[{"code":null,
"codeSystem":null,"name":"Tax or fee","kind":"unknown","payableAt":"prepaid","per":null,"jurisdiction":null,
"amount":{"amount":38.0,"currency":"EUR"},"approxAmount":{"amount":41.64,"currency":"USD"}}]},
"expiresAt":"2026-12-01T10:42:00Z","availability":"confirmed","checkInDate":"2026-12-10","checkOutDate":"2026-12-13",
"requestedRooms":[{"adults":2,"children":1,"childrenAges":[7]}],"terms":{"payment":{"type":"payNow",
"dueNow":{"amount":648.9,"currency":"USD"},"dueLater":null,"dueLaterAt":null},
"cancellation":{"type":"fullyRefundable","freeUntil":"2026-12-08T23:59:00Z",
"freeUntilLocalTime":"2026-12-09T00:59:00","timezone":"Europe/Madrid","schedule":[{"from":"2026-12-08T23:59:00Z",
"fromLocalTime":"2026-12-09T00:59:00","refund":{"amount":432.6,"currency":"USD"},"penalty":{"amount":216.3,
"currency":"USD"}}],"fee":null},"changes":{"allowed":true,"fee":null}},
"rate":{"rooms":[{"roomGroupKey":"4c8e3ae1d379f759bc372243","name":"Deluxe King","description":null,"beds":[],
"maxOccupancy":3,"sizeSquareMeters":null,"amenities":[],"imageUrls":[],"roomIndex":1,"adultsCount":2,
"childrenCount":1,"childrenAges":[7]}],"board":{"code":"BB","name":"Bed and breakfast","boardType":"breakfast"},
"priceIsBinding":true,"nightlyPrices":[{"date":"2026-12-10","amount":216.3,"currency":"USD"}],"rateTags":[],
"remarks":[],"rateConditions":[],"withTransfers":false,"rateRules":[],"amendabilityStatus":"yes","surcharges":[],
"minimumSellingPrice":null}},{"productType":"flight","publicId":"flight_guis4oiryxjtxfalagepbfdcpk","quantity":1,
"name":"LHR \u21c4 JFK \u00b7 12\u201319 Dec \u00b7 Economy Basic, Non-refundable","price":{"amount":1184.2,
"currency":"USD"},"taxes":{"treatment":"included","amount":{"amount":204.2,"currency":"USD"},"lines":[{"code":"US",
"codeSystem":"iata","name":"Transportation tax","kind":"tax","payableAt":"prepaid","per":null,"jurisdiction":"US",
"amount":{"amount":47.0,"currency":"USD"},"approxAmount":null}]},"expiresAt":"2026-12-01T10:30:00Z",
"quotedAt":"2026-12-01T10:05:12Z","hasSeatMaps":true,"terms":{"payment":{"type":"payNow","dueNow":{"amount":1184.2,
"currency":"USD"},"dueLater":null,"dueLaterAt":null},"cancellation":{"type":"nonRefundable","freeUntil":null,
"freeUntilLocalTime":null,"timezone":"Europe/London","schedule":[],"fee":null},"changes":{"allowed":true,"fee":null}},
"includedBaggage":[{"type":"carryOn","quantity":1,"weightKg":null,"description":null}],
"availableServices":[{"id":"flight_4224ss6vukylqlb2iphmrx3u52","type":"baggage","price":{"amount":30.0,
"currency":"USD"},"passengerIndexes":[1],"segments":[{"sliceIndex":1,"segmentIndex":1}],"maxQuantity":2,
"baggage":{"type":"checked","maximumWeightKg":23,"maximumLengthCm":null,"maximumHeightCm":null,
"maximumDepthCm":null}}],"fareOptions":[{"publicId":"flight_fi6uxnnuusc37omsl7vrs6z6lt",
"name":"Economy Basic, Non-refundable","price":{"amount":1184.2,"currency":"USD"},"taxes":{"treatment":"included",
"amount":{"amount":204.2,"currency":"USD"},"lines":[{"code":"US","codeSystem":"iata","name":"Transportation tax",
"kind":"tax","payableAt":"prepaid","per":null,"jurisdiction":"US","amount":{"amount":47.0,"currency":"USD"},
"approxAmount":null}]},"details":{"expiresAt":"2026-12-01T10:30:00Z","cabinClass":"economy",
"fare":{"fareBasisCode":"OXZ0RO","brandName":"Basic","fareFamily":"Economy Basic","validatingCarrier":"BA"},
"terms":{"payment":{"type":"payNow","dueNow":{"amount":1184.2,"currency":"USD"},"dueLater":null,"dueLaterAt":null},
"cancellation":{"type":"nonRefundable","freeUntil":null,"freeUntilLocalTime":null,"timezone":"Europe/London",
"schedule":[],"fee":null},"changes":{"allowed":true,"fee":null}},"includedBaggage":[{"type":"carryOn","quantity":1,
"weightKg":null,"description":null}],"passengerPricing":[{"passengerType":"adult","count":1,"pricePerPassenger":null,
"fareBasisCode":"OXZ0RO","bookingClass":"O"}],"offerPassengers":[{"passengerIndex":1,"type":"adult","age":null,
"fareType":null,"includedBaggage":[{"type":"carryOn","quantity":1,"weightKg":null,"description":null}],
"availableServices":[{"id":"flight_4224ss6vukylqlb2iphmrx3u52","type":"baggage","price":{"amount":30.0,
"currency":"USD"},"passengerIndexes":[1],"segments":[{"sliceIndex":1,"segmentIndex":1}],"maxQuantity":2,
"baggage":{"type":"checked","maximumWeightKg":23,"maximumLengthCm":null,"maximumHeightCm":null,
"maximumDepthCm":null}}]}],"availableServices":[{"id":"flight_4224ss6vukylqlb2iphmrx3u52","type":"baggage",
"price":{"amount":30.0,"currency":"USD"},"passengerIndexes":[1],"segments":[{"sliceIndex":1,"segmentIndex":1}],
"maxQuantity":2,"baggage":{"type":"checked","maximumWeightKg":23,"maximumLengthCm":null,"maximumHeightCm":null,
"maximumDepthCm":null}}],"identityDocumentsRequired":true,"supportedDocumentTypes":["passport"],
"segments":[{"sliceIndex":1,"segmentIndex":1,"cabinClass":"economy","baggage":[{"type":"carryOn","quantity":1,
"weightKg":null,"description":null}]}]}}]}],"total":{"amount":1885.6,"currency":"USD"},
"expiresAt":"2026-12-01T10:30:00Z"}`,
}

func readPayload[T any](t *testing.T, name string) T {
	t.Helper()
	var value T
	if err := json.Unmarshal([]byte(flatPayloads[name]), &value); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return value
}

// same reports every pair that differs; the pairs are (name, got, want).
func same(t *testing.T, pairs ...[3]any) {
	t.Helper()
	for _, p := range pairs {
		if p[1] != p[2] {
			t.Errorf("%v: got %v, want %v", p[0], p[1], p[2])
		}
	}
}

func TestAirtimeNumberLookupAndDeliveredValue(t *testing.T) {
	found := readPayload[AirtimeSearchResult](t, "airtimeSearch")
	quoted := readPayload[AirtimeQuote](t, "airtimeQuote").Items[0]
	same(t,
		[3]any{"phoneLookup.status", string(found.PhoneLookup.Status), "detected"},
		[3]any{"phoneLookup.carrierName", found.PhoneLookup.GetCarrierName(), "Bakcell LLC"},
		[3]any{"offer details", found.Items[0].Offers[0].Details.GetReadMore(), "Credit arrives within minutes."},
		[3]any{"delivered amount", float64(quoted.EstimatedDeliveredValue.Amount), float64(40)},
		[3]any{"delivered currency", quoted.EstimatedDeliveredValue.Currency, "AZN"},
	)
}

func TestEsimDestinationFieldsAreOnTheItem(t *testing.T) {
	destination := readPayload[EsimSearchResult](t, "esimSearch").Items[0]
	same(t,
		[3]any{"destinationCode", destination.GetDestinationCode(), "EU"},
		[3]any{"region", destination.GetRegion(), "Europe"},
		[3]any{"coverage", destination.Coverage[0].CountryCode, "DE"},
		[3]any{"offer details", int(destination.Offers[0].Details.GetValidityDays()), 7},
	)
}

func TestGiftCardAndUtilityFieldsAreBesideTheCommonKeys(t *testing.T) {
	card := readPayload[GiftCardSearchResult](t, "giftCardSearch").Items[0]
	biller := readPayload[UtilitySearchResult](t, "utilitySearch").Items[0]
	quoted := readPayload[UtilityQuote](t, "utilityQuote").Items[0]
	same(t,
		[3]any{"isRedeemableGlobally", card.IsRedeemableGlobally, false},
		[3]any{"redemptionInstructions", card.GetRedemptionInstructions(), "Go to amazon.com/redeem and enter the code."},
		[3]any{"serviceType", string(biller.ServiceType), "postpaid"},
		[3]any{"deliveryType", string(biller.DeliveryType), "directCredit"},
		[3]any{"quote deliveryType", string(quoted.DeliveryType), "directCredit"},
	)
}

func TestHotelFactsAndQuotedStay(t *testing.T) {
	found := readPayload[HotelSearchResult](t, "hotelSearch")
	hotel := found.Items[0]
	catalog := readPayload[CatalogHotel](t, "catalogHotel")
	stay := readPayload[HotelQuote](t, "hotelQuote").Items[0]
	filters := found.AvailableFilters.Get()
	if hotel.Address.Get() == nil || hotel.Coordinates.Get() == nil || catalog.Address.Get() == nil || filters == nil {
		t.Fatalf("a hotel fact is missing: %+v %+v", hotel, catalog)
	}
	same(t,
		[3]any{"starRating", float64(hotel.GetStarRating()), float64(5)},
		[3]any{"city", hotel.Address.Get().GetCity(), "Barcelona"},
		[3]any{"latitude", float64(hotel.Coordinates.Get().Latitude), 41.386},
		[3]any{"images", len(hotel.ImageUrls), 1},
		[3]any{"amenities", hotel.AmenityGroups[0].Codes[0], "FreeWifi"},
		[3]any{"offer details", hotel.Offers[0].Details.Board.GetCode(), "BB"},
		[3]any{"availableFilters", filters.Amenities[0].Code, "FreeWifi"},
		[3]any{"catalog starRating", float64(catalog.GetStarRating()), float64(5)},
		[3]any{"catalog postalCode", catalog.Address.Get().GetPostalCode(), "08005"},
		[3]any{"availability", string(stay.Availability), "confirmed"},
		[3]any{"expiresAt", stay.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"), "2026-12-01T10:42:00Z"},
		[3]any{"checkInDate", stay.CheckInDate, "2026-12-10"},
		[3]any{"checkOutDate", stay.CheckOutDate, "2026-12-13"},
		[3]any{"childrenAges", int(stay.RequestedRooms[0].ChildrenAges[0]), 7},
		[3]any{"payment type", string(stay.Terms.Payment.Type), "payNow"},
		[3]any{"board", stay.Rate.Board.GetCode(), "BB"},
	)
}

func TestFlightItineraryPlaceAndQuoteItemFields(t *testing.T) {
	itinerary := readPayload[FlightSearchResult](t, "flightSearch").Items[0]
	place := readPayload[FlightPlaceList](t, "flightPlaces").Items[0]
	fare := readPayload[FlightQuote](t, "flightQuote").Items[0]
	same(t,
		[3]any{"owner", itinerary.Owner.GetIataCode(), "BA"},
		[3]any{"flight number", itinerary.Slices[0].Segments[0].GetMarketingFlightNumber(), "117"},
		[3]any{"marketing carrier", itinerary.MarketingCarriers[0].GetIataCode(), "BA"},
		[3]any{"operating carrier", itinerary.OperatingCarriers[0].GetIataCode(), "BA"},
		[3]any{"duration", int(itinerary.TotalDurationMinutes), 885},
		[3]any{"stops", int(itinerary.TotalStops), 0},
		[3]any{"emissions", float64(itinerary.GetTotalEmissionsKg()), float64(1236)},
		[3]any{"offer details", string(itinerary.Offers[0].Details.CabinClass), "economy"},
		[3]any{"iataCode", place.IataCode, "LON"},
		[3]any{"cityName", place.CityName, "London"},
		[3]any{"cityIataCode", place.GetCityIataCode(), "LON"},
		[3]any{"timezone", place.GetTimezone(), "Europe/London"},
		[3]any{"airport", place.Airports[0].IataCode, "LGW"},
		[3]any{"expiresAt", fare.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"), "2026-12-01T10:30:00Z"},
		[3]any{"quotedAt", fare.QuotedAt.UTC().Format("2006-01-02T15:04:05Z"), "2026-12-01T10:05:12Z"},
		[3]any{"hasSeatMaps", fare.HasSeatMaps, true},
		[3]any{"cancellation", string(fare.Terms.Cancellation.Type), "nonRefundable"},
		[3]any{"baggage", string(fare.IncludedBaggage[0].Type), "carryOn"},
		[3]any{"maxQuantity", int(fare.AvailableServices[0].GetMaxQuantity()), 2},
		[3]any{"fare option details", string(fare.FareOptions[0].Details.CabinClass), "economy"},
	)
}

func TestMixedQuoteVariantsCarryTheSameFields(t *testing.T) {
	items := readPayload[MixedQuote](t, "mixedQuote").Items
	airtime, utility := items[0].AirtimeMixedQuoteItem, items[1].UtilityMixedQuoteItem
	hotel, flight := items[2].HotelMixedQuoteItem, items[3].FlightMixedQuoteItem
	if airtime == nil || utility == nil || hotel == nil || flight == nil {
		t.Fatalf("a variant is missing: %+v", items)
	}
	same(t,
		[3]any{"delivered amount", float64(airtime.EstimatedDeliveredValue.Amount), float64(40)},
		[3]any{"deliveryType", string(utility.DeliveryType), "directCredit"},
		[3]any{"availability", string(hotel.Availability), "confirmed"},
		[3]any{"board", hotel.Rate.Board.GetCode(), "BB"},
		[3]any{"hasSeatMaps", flight.HasSeatMaps, true},
		[3]any{"fare option details", string(flight.FareOptions[0].Details.CabinClass), "economy"},
	)
}
