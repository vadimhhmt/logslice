package flatten_test

import (
	"testing"

	"logslice/internal/flatten"
	"logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestFlattener_FlatEntryUnchanged(t *testing.T) {
	f := flatten.New(".", 0)
	input := makeEntry(map[string]interface{}{"level": "info", "msg": "hello"})
	out := f.Process(input)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatalf("unexpected output: %v", out)
	}
}

func TestFlattener_NestedMapFlattened(t *testing.T) {
	f := flatten.New(".", 0)
	input := makeEntry(map[string]interface{}{
		"user": map[string]interface{}{"id": "42", "name": "alice"},
	})
	out := f.Process(input)
	if _, exists := out["user"]; exists {
		t.Fatal("nested key 'user' should have been removed")
	}
	if out["user.id"] != "42" {
		t.Errorf("expected user.id=42, got %v", out["user.id"])
	}
	if out["user.name"] != "alice" {
		t.Errorf("expected user.name=alice, got %v", out["user.name"])
	}
}

func TestFlattener_DeeplyNested(t *testing.T) {
	f := flatten.New(".", 0)
	input := makeEntry(map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	})
	out := f.Process(input)
	if out["a.b.c"] != "deep" {
		t.Errorf("expected a.b.c=deep, got %v", out["a.b.c"])
	}
}

func TestFlattener_MaxDepthLimitsNesting(t *testing.T) {
	f := flatten.New(".", 1)
	input := makeEntry(map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	})
	out := f.Process(input)
	// at depth 1 the inner map should be stored as a string, not further flattened
	if _, exists := out["a.b.c"]; exists {
		t.Error("a.b.c should not exist when maxDepth=1")
	}
	if out["a.b"] == nil {
		t.Error("a.b should exist as a stringified value")
	}
}

func TestFlattener_OriginalUnmodified(t *testing.T) {
	f := flatten.New(".", 0)
	nested := map[string]interface{}{"x": 1}
	input := makeEntry(map[string]interface{}{"meta": nested})
	f.Process(input)
	if _, exists := input["meta"]; !exists {
		t.Error("original entry should not be mutated")
	}
}

func TestFlattener_CustomSeparator(t *testing.T) {
	f := flatten.New("_", 0)
	input := makeEntry(map[string]interface{}{
		"http": map[string]interface{}{"status": 200},
	})
	out := f.Process(input)
	if out["http_status"] != 200 {
		t.Errorf("expected http_status=200, got %v", out["http_status"])
	}
}
