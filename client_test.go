package raven

import (
	"testing"
	"time"
)

func TestPacketJSON(t *testing.T) {
	packet := &Packet{
		Project:    "1",
		EventID:    "2",
		Message:    "test",
		Timestamp:  Timestamp(time.Date(2000, 01, 01, 0, 0, 0, 0, time.UTC)),
		Level:      ERROR,
		Interfaces: []Interface{&Message{Message: "foo"}},
	}

	expected := `{"message":"test","event_id":"2","project":"1","timestamp":"2000-01-01T00:00:00","level":40,"sentry.interfaces.Message":{"message":"foo"}}`
	actual := string(packet.JSON())

	if actual != expected {
		t.Errorf("incorrect json; got %s, want %s", actual, expected)
	}
}

func TestPacketInit(t *testing.T) {
	packet := &Packet{Message: "a", Tags: map[string]string{"foo": "bar"}}
	packet.Init("foo", map[string]string{"foo": "foo", "baz": "buzz"})

	if packet.Project != "foo" {
		t.Error("incorrect Project:", packet.Project)
	}
	if packet.Level != ERROR {
		t.Errorf("incorrect Level: got %d, want %d", packet.Level, ERROR)
	}
	if time.Time(packet.Timestamp).IsZero() {
		t.Error("Timestamp is zero")
	}
	if len(packet.EventID) != 32 {
		t.Error("incorrect EventID:", packet.EventID)
	}
	if len(packet.Tags) != 2 || packet.Tags["foo"] != "bar" || packet.Tags["baz"] != "buzz" {
		t.Errorf("incorrect Tags: %#v", packet.Tags)
	}
}

func TestSetDSN(t *testing.T) {
	client := &Client{}
	client.SetDSN("https://u:p@example.com/sentry/1")

	if client.url != "https://example.com/sentry/api/1/store/" {
		t.Error("incorrect url:", client.url)
	}
	if client.projectID != "1" {
		t.Error("incorrect projectID:", client.projectID)
	}
	if client.authHeader != "Sentry sentry_version=3, sentry_key=u, sentry_secret=p" {
		t.Error("incorrect authHeader:", client.authHeader)
	}
}
