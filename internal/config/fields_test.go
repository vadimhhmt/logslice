package config

import (
	"reflect"
	"testing"
)

func TestFieldList_Empty(t *testing.T) {
	if got := FieldList(""); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestFieldList_Single(t *testing.T) {
	got := FieldList("level")
	want := []string{"level"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFieldList_Multiple(t *testing.T) {
	got := FieldList("level,msg,service")
	want := []string{"level", "msg", "service"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFieldList_Deduplication(t *testing.T) {
	got := FieldList("level,msg,level")
	if len(got) != 2 {
		t.Errorf("expected 2 unique fields, got %d: %v", len(got), got)
	}
}

func TestFieldList_Whitespace(t *testing.T) {
	got := FieldList(" level , msg ")
	want := []string{"level", "msg"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestPatternPairs_Empty(t *testing.T) {
	got := PatternPairs("")
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestPatternPairs_Single(t *testing.T) {
	got := PatternPairs("level=error")
	if got["level"] != "error" {
		t.Errorf("expected level=error, got %v", got)
	}
}

func TestPatternPairs_Multiple(t *testing.T) {
	got := PatternPairs("level=error,service=api")
	if got["level"] != "error" || got["service"] != "api" {
		t.Errorf("unexpected pairs: %v", got)
	}
}

func TestPatternPairs_MissingEquals(t *testing.T) {
	got := PatternPairs("levelonly")
	if len(got) != 0 {
		t.Errorf("expected empty map for missing '=', got %v", got)
	}
}

func TestPatternPairs_ValueWithEquals(t *testing.T) {
	// value itself contains '=' — only first '=' is the separator
	got := PatternPairs("msg=hello=world")
	if got["msg"] != "hello=world" {
		t.Errorf("expected msg=hello=world, got %v", got)
	}
}
