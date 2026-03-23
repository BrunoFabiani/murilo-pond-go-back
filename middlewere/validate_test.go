package middlewere

import (
	"testing"
	"time"
)

func TestValidateTelemetrySuccessAnalog(t *testing.T) {
	msg := Telemetry{
		DeviceID:      "device-1",
		Timestamp:     time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
		SensorType:    "temperature",
		ReadingNature: "analog",
		Value:         23.5,
	}

	if err := validateTelemetry(&msg); err != nil {
		t.Fatalf("expected valid telemetry, got error: %v", err)
	}
	if msg.DeviceID != "device-1" {
		t.Fatalf("expected trimmed device_id to stay 'device-1', got %q", msg.DeviceID)
	}
}

func TestValidateTelemetrySuccessDiscrete(t *testing.T) {
	msg := Telemetry{
		DeviceID:      "device-2",
		Timestamp:     time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
		SensorType:    "presence",
		ReadingNature: "discrete",
		Value:         true,
	}

	if err := validateTelemetry(&msg); err != nil {
		t.Fatalf("expected valid discrete telemetry, got error: %v", err)
	}
}

func TestValidateTelemetryMissingDeviceID(t *testing.T) {
	msg := Telemetry{
		Timestamp:     time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
		SensorType:    "temperature",
		ReadingNature: "analog",
		Value:         23.5,
	}

	if err := validateTelemetry(&msg); err == nil {
		t.Fatal("expected validation error for missing device_id")
	}
}

func TestValidateTelemetryInvalidReadingNature(t *testing.T) {
	msg := Telemetry{
		DeviceID:      "device-3",
		Timestamp:     time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
		SensorType:    "temperature",
		ReadingNature: "invalid",
		Value:         23.5,
	}

	if err := validateTelemetry(&msg); err == nil {
		t.Fatal("expected validation error for invalid reading_nature")
	}
}

func TestValidateTelemetryAnalogRequiresNumericValue(t *testing.T) {
	msg := Telemetry{
		DeviceID:      "device-4",
		Timestamp:     time.Date(2026, time.March, 22, 10, 0, 0, 0, time.UTC),
		SensorType:    "temperature",
		ReadingNature: "analog",
		Value:         "hot",
	}

	if err := validateTelemetry(&msg); err == nil {
		t.Fatal("expected validation error for non-numeric analog value")
	}
}
