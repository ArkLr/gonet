package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStatRoundTrip(t *testing.T) {
	want := StatsMsg{
		Time:       time.Date(2026, 11, 17, 9, 0, 0, 0, time.UTC),
		IntervalMs: 1000,
		BytesIn:    4096,
		BytesOut:   512,
		PacketsIn:  12,
		PacketsOut: 9,
		Dropped:    3,
	}

	//Encode the message as the agent will
	data, err := json.Marshal(NewStats(want))
	if err != nil {
		t.Fatalf("encoding failed with marshal: %v", err)
	}
	t.Logf("wire: %s", data)

	//Decode the message as the agent will
	var got Message
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("decoding failed with unmarshal: %v", err)
	}

	//Validate the message recived
	if err := got.Validate(); err != nil {
		t.Fatalf("message validation error: %v", err)
	}

	if *got.Stats != want {
		t.Errorf("stats changed in round trip:\n got  %+v\n want %+v", *got.Stats, want)
	}
}

func TestValidateRejectsMissingPayload(t *testing.T) {
	bad := Message{
		Version: version,
		Type:    TypeStats,
	}

	if err := bad.Validate(); err == nil {
		t.Fatal("expected an error for a stats message with no payload, got nil")
	}
}
