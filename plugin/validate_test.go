package main

import (
	"context"
	"testing"

	"github.com/MrWong99/zhi/pkg/zhiplugin/config"
)

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		blocking bool
	}{
		{"empty string blocks", "", true},
		{"non-empty passes", "hello", false},
		{"non-string blocks", 42, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := validateRequired(config.Value{Val: tt.val}, nil)
			if err != nil {
				t.Fatal(err)
			}
			hasBlocking := len(results) > 0 && results[0].Severity == config.Blocking
			if hasBlocking != tt.blocking {
				t.Errorf("blocking = %v, want %v (results: %v)", hasBlocking, tt.blocking, results)
			}
		})
	}
}

func TestValidateAbsolutePath(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		blocking bool
	}{
		{"absolute path passes", "/srv/data", false},
		{"relative path blocks", "srv/data", true},
		{"empty string blocks", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := validateAbsolutePath(config.Value{Val: tt.val}, nil)
			if err != nil {
				t.Fatal(err)
			}
			hasBlocking := len(results) > 0 && results[0].Severity == config.Blocking
			if hasBlocking != tt.blocking {
				t.Errorf("blocking = %v, want %v", hasBlocking, tt.blocking)
			}
		})
	}
}

func TestValidateOptionalAbsPath(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		blocking bool
	}{
		{"empty passes", "", false},
		{"absolute passes", "/mnt/media", false},
		{"relative blocks", "mnt/media", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := validateOptionalAbsPath(config.Value{Val: tt.val}, nil)
			if err != nil {
				t.Fatal(err)
			}
			hasBlocking := len(results) > 0 && results[0].Severity == config.Blocking
			if hasBlocking != tt.blocking {
				t.Errorf("blocking = %v, want %v", hasBlocking, tt.blocking)
			}
		})
	}
}

func TestValidatePiholeDNSPort(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		severity config.Severity
		hasResult bool
	}{
		{"port 53 blocks", 53, config.Blocking, true},
		{"port 5353 passes", 5353, 0, false},
		{"non-number blocks", "abc", config.Blocking, true},
		{"float port passes", 8053.0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := validatePiholeDNSPort(config.Value{Val: tt.val}, nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.hasResult {
				if len(results) == 0 {
					t.Fatal("expected validation result, got none")
				}
				if results[0].Severity != tt.severity {
					t.Errorf("severity = %v, want %v", results[0].Severity, tt.severity)
				}
			} else if len(results) > 0 {
				t.Errorf("expected no results, got %v", results)
			}
		})
	}
}

func TestValidatorsMapOnlyReferencesKnownPaths(t *testing.T) {
	known := make(map[string]bool, len(valueDefs))
	for _, d := range valueDefs {
		known[d.Path] = true
	}
	for path := range validators {
		if !known[path] {
			t.Errorf("validator registered for unknown path: %s", path)
		}
	}
}

func TestPluginValidateDispatch(t *testing.T) {
	p := newHomeserverPlugin()

	// Build a tree reader from the plugin's own defaults.
	tree := config.NewTree()
	paths, _ := p.List(context.Background())
	for _, path := range paths {
		v, ok, _ := p.Get(context.Background(), path)
		if ok {
			tree.Set(path, &v)
		}
	}

	// core/domain defaults to "" which is required -- should block.
	results, err := p.Validate(context.Background(), "core/domain", tree)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected blocking result for empty core/domain")
	}
	if results[0].Severity != config.Blocking {
		t.Errorf("severity = %v, want Blocking", results[0].Severity)
	}

	// core/timezone has no validator -- should return nil.
	results, err = p.Validate(context.Background(), "core/timezone", tree)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for core/timezone, got %v", results)
	}
}
