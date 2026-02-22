package main

import "github.com/MrWong99/zhi/pkg/zhiplugin/config"

// ValueDef defines a configuration value with its default and metadata.
// This reduces the boilerplate of repeating the same metadata label keys
// across all config values.
type ValueDef struct {
	Path        string // slash-delimited config path, e.g. "core/timezone"
	Default     any    // default value
	Section     string // ui.section
	DisplayName string // ui.displayName
	Description string // core.description
	Type        string // core.type (string, int, bool)

	// Optional fields -- zero values mean "not set"
	Placeholder string   // ui.placeholder
	Password    bool     // ui.password
	Required    bool     // config.required
	SelectFrom  []string // ui.enum (dropdown selection)
}

// ToValue converts a ValueDef to a config.Value with the standard
// metadata labels populated.
func (d *ValueDef) ToValue() *config.Value {
	md := map[string]any{
		"ui.section":       d.Section,
		"ui.displayName":   d.DisplayName,
		"core.description": d.Description,
		"core.type":        d.Type,
	}
	if d.Placeholder != "" {
		md["ui.placeholder"] = d.Placeholder
	}
	if d.Password {
		md["ui.password"] = true
	}
	if d.Required {
		md["config.required"] = true
	}
	if len(d.SelectFrom) > 0 {
		md["ui.enum"] = d.SelectFrom
	}
	return &config.Value{
		Val:      d.Default,
		Metadata: md,
	}
}
