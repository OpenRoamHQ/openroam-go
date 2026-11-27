package openroam

import (
	"encoding/json"
	"os"
)

// testdata/samples.json (generated from the API description) holds the smallest valid value of every type: its required
// properties only. The SDK refuses an answer without a required property, so a test that needs an answer of a type starts
// from its sample and sets only what it checks.
var samples = func() map[string]json.RawMessage {
	data, err := os.ReadFile("testdata/samples.json")
	if err != nil {
		panic(err)
	}
	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		panic(err)
	}
	return all
}()

// sample is the JSON of the sample of a type; it panics on a type the API description does not have.
func sample(name string) string {
	raw, ok := samples[name]
	if !ok {
		panic("no sample for " + name)
	}
	return string(raw)
}

// sampleWith is the sample of a type as an object, with fields set (replaced) on top.
func sampleWith(name string, fields map[string]any) map[string]any {
	var object map[string]any
	if err := json.Unmarshal([]byte(sample(name)), &object); err != nil {
		panic(err)
	}
	for key, value := range fields {
		object[key] = value
	}
	return object
}

// sampleJSON is sampleWith encoded as JSON.
func sampleJSON(name string, fields map[string]any) []byte {
	data, err := json.Marshal(sampleWith(name, fields))
	if err != nil {
		panic(err)
	}
	return data
}
