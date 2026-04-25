package transform_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/transform"
)

func makeEntry(fields map[string]interface{}) parser.LogEntry {
	return parser.LogEntry{
		Timestamp: time.Now(),
		Raw:       `{"level":"info"}`,
		Fields:    fields,
	}
}

func TestTransformer_RenameField(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{FromField: "msg", ToField: "message"},
	})
	e := makeEntry(map[string]interface{}{"msg": "hello", "level": "info"})
	out := tr.Apply(e)
	if _, ok := out.Fields["msg"]; ok {
		t.Error("old field 'msg' should be removed")
	}
	if out.Fields["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", out.Fields["message"])
	}
}

func TestTransformer_ValueMap(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{FromField: "level", ValueMap: map[string]string{"warn": "warning"}},
	})
	e := makeEntry(map[string]interface{}{"level": "warn"})
	out := tr.Apply(e)
	if out.Fields["level"] != "warning" {
		t.Errorf("expected level=warning, got %v", out.Fields["level"])
	}
}

func TestTransformer_NoMatchLeaveIntact(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{FromField: "nonexistent", ToField: "other"},
	})
	e := makeEntry(map[string]interface{}{"level": "info"})
	out := tr.Apply(e)
	if _, ok := out.Fields["other"]; ok {
		t.Error("unexpected field 'other' should not exist")
	}
	if out.Fields["level"] != "info" {
		t.Error("unrelated field should remain unchanged")
	}
}

func TestTransformer_OriginalUnmodified(t *testing.T) {
	tr := transform.New([]transform.Rule{
		{FromField: "msg", ToField: "message"},
	})
	e := makeEntry(map[string]interface{}{"msg": "hi"})
	_ = tr.Apply(e)
	if _, ok := e.Fields["msg"]; !ok {
		t.Error("original entry should not be modified")
	}
}

func TestTransformer_EmptyRules(t *testing.T) {
	tr := transform.New(nil)
	e := makeEntry(map[string]interface{}{"level": "debug", "msg": "test"})
	out := tr.Apply(e)
	if len(out.Fields) != len(e.Fields) {
		t.Errorf("expected %d fields, got %d", len(e.Fields), len(out.Fields))
	}
}
