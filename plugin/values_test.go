package main

import (
	"context"
	"testing"

	"github.com/MrWong99/zhi/pkg/zhiplugin/config"
)

func TestAllPathsAreUnique(t *testing.T) {
	seen := make(map[string]bool, len(valueDefs))
	for _, d := range valueDefs {
		if seen[d.Path] {
			t.Errorf("duplicate path: %s", d.Path)
		}
		seen[d.Path] = true
	}
}

func TestAllPathsAreValid(t *testing.T) {
	for _, d := range valueDefs {
		if err := config.ValidatePath(d.Path); err != nil {
			t.Errorf("invalid path %q: %v", d.Path, err)
		}
	}
}

func TestAllDefsHaveRequiredMetadata(t *testing.T) {
	for _, d := range valueDefs {
		if d.Section == "" {
			t.Errorf("path %q missing Section", d.Path)
		}
		if d.DisplayName == "" {
			t.Errorf("path %q missing DisplayName", d.Path)
		}
		if d.Description == "" {
			t.Errorf("path %q missing Description", d.Path)
		}
		if d.Type == "" {
			t.Errorf("path %q missing Type", d.Path)
		}
	}
}

func TestToValueMetadata(t *testing.T) {
	d := ValueDef{
		Path: "test/value", Default: "hello",
		Section: "General", DisplayName: "Test",
		Description: "A test value", Type: "string",
		Placeholder: "enter value", Password: true, Required: true,
		SelectFrom: []string{"a", "b"},
	}
	v := d.ToValue()

	if v.Val != "hello" {
		t.Errorf("Val = %v, want hello", v.Val)
	}
	checks := map[string]any{
		"ui.section":       "General",
		"ui.displayName":   "Test",
		"core.description": "A test value",
		"core.type":        "string",
		"ui.placeholder":   "enter value",
		"ui.password":      true,
		"config.required":  true,
	}
	for k, want := range checks {
		got, ok := v.Metadata[k]
		if !ok {
			t.Errorf("metadata key %q missing", k)
			continue
		}
		if got != want {
			t.Errorf("metadata[%q] = %v, want %v", k, got, want)
		}
	}
	enumVal, ok := v.Metadata["ui.enum"]
	if !ok {
		t.Fatal("metadata key ui.enum missing")
	}
	enum, ok := enumVal.([]string)
	if !ok {
		t.Fatalf("ui.enum is %T, want []string", enumVal)
	}
	if len(enum) != 2 || enum[0] != "a" || enum[1] != "b" {
		t.Errorf("ui.enum = %v, want [a b]", enum)
	}
}

func TestToValueOmitsZeroOptionals(t *testing.T) {
	d := ValueDef{
		Path: "test/minimal", Default: 42,
		Section: "S", DisplayName: "D",
		Description: "Desc", Type: "int",
	}
	v := d.ToValue()

	for _, key := range []string{"ui.placeholder", "ui.password", "config.required", "ui.enum"} {
		if _, ok := v.Metadata[key]; ok {
			t.Errorf("metadata key %q should not be set for zero-value optionals", key)
		}
	}
}

func TestPluginListReturnsAllPaths(t *testing.T) {
	p := newHomeserverPlugin()
	paths, err := p.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != len(valueDefs) {
		t.Errorf("List returned %d paths, want %d", len(paths), len(valueDefs))
	}
}

func TestPluginGetReturnsValues(t *testing.T) {
	p := newHomeserverPlugin()
	for _, d := range valueDefs {
		v, ok, err := p.Get(context.Background(), d.Path)
		if err != nil {
			t.Errorf("Get(%q): %v", d.Path, err)
			continue
		}
		if !ok {
			t.Errorf("Get(%q): not found", d.Path)
			continue
		}
		if v.Metadata == nil {
			t.Errorf("Get(%q): metadata is nil", d.Path)
		}
	}
}

func TestPluginGetMissingPath(t *testing.T) {
	p := newHomeserverPlugin()
	_, ok, err := p.Get(context.Background(), "nonexistent/path")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("Get for nonexistent path returned ok=true")
	}
}

func TestPluginSetAndGet(t *testing.T) {
	p := newHomeserverPlugin()
	err := p.Set(context.Background(), "core/timezone", config.Value{Val: "UTC"})
	if err != nil {
		t.Fatal(err)
	}
	v, ok, err := p.Get(context.Background(), "core/timezone")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("path not found after Set")
	}
	if v.Val != "UTC" {
		t.Errorf("Val = %v, want UTC", v.Val)
	}
}

func TestPluginSetInvalidPath(t *testing.T) {
	p := newHomeserverPlugin()
	err := p.Set(context.Background(), "INVALID", config.Value{Val: "x"})
	if err == nil {
		t.Error("Set with invalid path should return error")
	}
}
