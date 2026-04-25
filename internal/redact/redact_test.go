package redact_test

import (
	"testing"

	"logslice/internal/redact"
)

func TestRedactor_DefaultSensitiveKeys(t *testing.T) {
	r := redact.New()
	fields := map[string]any{
		"user":     "alice",
		"password": "s3cr3t",
		"token":    "eyJhb...",
	}
	out := r.Apply(fields)
	if out["user"] != "alice" {
		t.Errorf("expected user to be unchanged, got %v", out["user"])
	}
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected password to be redacted, got %v", out["password"])
	}
	if out["token"] != "[REDACTED]" {
		t.Errorf("expected token to be redacted, got %v", out["token"])
	}
}

func TestRedactor_CaseInsensitiveKey(t *testing.T) {
	r := redact.New()
	fields := map[string]any{"Password": "hunter2", "API_KEY": "abc123"}
	out := r.Apply(fields)
	if out["Password"] != "[REDACTED]" {
		t.Errorf("expected Password redacted, got %v", out["Password"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY redacted, got %v", out["API_KEY"])
	}
}

func TestRedactor_AddKey(t *testing.T) {
	r := redact.New()
	r.AddKey("ssn")
	fields := map[string]any{"ssn": "123-45-6789", "name": "Bob"}
	out := r.Apply(fields)
	if out["ssn"] != "[REDACTED]" {
		t.Errorf("expected ssn redacted, got %v", out["ssn"])
	}
	if out["name"] != "Bob" {
		t.Errorf("expected name unchanged, got %v", out["name"])
	}
}

func TestRedactor_AddPattern(t *testing.T) {
	r := redact.New()
	if err := r.AddPattern(`^Bearer\s+\S+`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := map[string]any{
		"header": "Bearer eyJhbGciOiJSUzI1",
		"other":  "plain text",
	}
	out := r.Apply(fields)
	if out["header"] != "[REDACTED]" {
		t.Errorf("expected header redacted, got %v", out["header"])
	}
	if out["other"] != "plain text" {
		t.Errorf("expected other unchanged, got %v", out["other"])
	}
}

func TestRedactor_InvalidPattern(t *testing.T) {
	r := redact.New()
	if err := r.AddPattern(`[invalid`); err == nil {
		t.Error("expected error for invalid pattern, got nil")
	}
}

func TestRedactor_NoSensitiveFields(t *testing.T) {
	r := redact.New()
	fields := map[string]any{"level": "info", "msg": "hello"}
	out := r.Apply(fields)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Errorf("expected fields unchanged, got %v", out)
	}
}

func TestRedactor_NonStringValueNotMatchedByPattern(t *testing.T) {
	r := redact.New()
	_ = r.AddPattern(`\d+`)
	fields := map[string]any{"count": 42}
	out := r.Apply(fields)
	if out["count"] != 42 {
		t.Errorf("expected numeric count unchanged, got %v", out["count"])
	}
}
