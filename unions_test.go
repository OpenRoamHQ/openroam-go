package openroam

import (
	"encoding/json"
	"strings"
	"testing"
)

// Responses and requests carry a union of one variant per product, chosen by the `productType` property.
func TestOrderItemIsReadAsItsOwnProductVariant(t *testing.T) {
	for productType, check := range map[string]func(PartnerOrderItem) bool{
		"airtime":  func(i PartnerOrderItem) bool { return i.PartnerOrderItemAirtime != nil },
		"esim":     func(i PartnerOrderItem) bool { return i.PartnerOrderItemEsim != nil },
		"flight":   func(i PartnerOrderItem) bool { return i.PartnerOrderItemFlight != nil },
		"giftCard": func(i PartnerOrderItem) bool { return i.PartnerOrderItemGiftCard != nil },
		"hotel":    func(i PartnerOrderItem) bool { return i.PartnerOrderItemHotel != nil },
		"utility":  func(i PartnerOrderItem) bool { return i.PartnerOrderItemUtility != nil },
	} {
		var item PartnerOrderItem
		if err := json.Unmarshal([]byte(`{"id":"11111111-1111-1111-1111-111111111111","productType":"`+productType+`"}`), &item); err != nil {
			t.Fatalf("%s: %v", productType, err)
		}
		if !check(item) {
			t.Errorf("%s: the wrong variant was set: %+v", productType, item)
		}
	}
}

func TestEsimWebhookItemKeepsItsDeliveryDetails(t *testing.T) {
	var item OrderFulfillmentUpdatedWebhookItem
	body := `{"id":"11111111-1111-1111-1111-111111111111","productType":"esim","fulfillment":{"status":"fulfilled","result":{"details":{"iccid":"8900000000000000001","mode":"instant"}}}}`
	if err := json.Unmarshal([]byte(body), &item); err != nil {
		t.Fatal(err)
	}
	if item.OrderFulfillmentUpdatedWebhookItemEsim == nil {
		t.Fatalf("not the eSIM variant: %+v", item)
	}
	out, err := json.Marshal(item)
	if err != nil || !strings.Contains(string(out), "8900000000000000001") {
		t.Fatalf("the ICCID did not survive a round trip: %s %v", out, err)
	}
}

func TestPendingWebhookItemHasNoResult(t *testing.T) {
	var item OrderFulfillmentUpdatedWebhookItem
	body := `{"id":"11111111-1111-1111-1111-111111111111","productType":"airtime","fulfillment":{"status":"pending","result":null}}`
	if err := json.Unmarshal([]byte(body), &item); err != nil {
		t.Fatal(err)
	}
	if item.OrderFulfillmentUpdatedWebhookItemAirtime == nil {
		t.Fatalf("not the airtime variant: %+v", item)
	}
}

func TestAnUnknownProductTypeIsAnErrorNotASilentGap(t *testing.T) {
	var item PartnerOrderItem
	err := json.Unmarshal([]byte(`{"id":"1","productType":"spaceship"}`), &item)
	if err == nil || !strings.Contains(err.Error(), "spaceship") {
		t.Fatalf("want an error naming the unknown productType, got %v", err)
	}
}

func TestOrderCommandItemsKeepTheirVariant(t *testing.T) {
	var command CreateB2BOrderCommand
	if err := json.Unmarshal([]byte(`{"items":[{"productType":"esim","publicId":"esim_abc","quantity":1}]}`), &command); err != nil {
		t.Fatal(err)
	}
	if len(command.Items) != 1 || command.Items[0].CreateB2BOrderItemCommandEsim == nil {
		t.Fatalf("the esim item was lost: %+v", command.Items)
	}
}

// An offer's price is a fixed amount or a range, and its value is discriminated by `type`.
func TestOfferValueAndPriceAreReadByTheirShape(t *testing.T) {
	var offer AirtimeOffer
	body := `{"publicId":"airtime_abc","name":"MTN 1000 NGN","value":{"type":"fixed","amount":1000,"currency":"NGN"},"price":{"amount":0.68,"currency":"USD"}}`
	if err := json.Unmarshal([]byte(body), &offer); err != nil {
		t.Fatal(err)
	}
	if offer.Value == nil || offer.Value.OfferFixedValue == nil {
		t.Fatalf("the value is not the fixed variant: %+v", offer.Value)
	}
	if offer.Price == nil || offer.Price.OfferFixedPrice == nil {
		t.Fatalf("the price is not the fixed variant: %+v", offer.Price)
	}

	var ranged AirtimeOffer
	body = `{"publicId":"airtime_abd","name":"MTN range","value":{"type":"range","min":50,"max":50000,"step":1,"currency":"NGN"},"price":{"min":0.04,"max":33.5,"currency":"USD"}}`
	if err := json.Unmarshal([]byte(body), &ranged); err != nil {
		t.Fatal(err)
	}
	if ranged.Value.OfferRangeValue == nil || ranged.Price.OfferRangePrice == nil {
		t.Fatalf("the range variants were not read: %+v %+v", ranged.Value, ranged.Price)
	}
}
