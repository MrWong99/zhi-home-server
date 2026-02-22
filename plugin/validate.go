package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MrWong99/zhi/pkg/zhiplugin/config"
)

// validatorFunc validates a single value, optionally using the full tree
// for cross-value checks.
type validatorFunc func(v config.Value, tree config.TreeReader) ([]config.ValidationResult, error)

// validators maps config paths to their validation functions.
// Only paths that need validation are listed -- unlisted paths are always valid.
var validators = map[string]validatorFunc{
	"core/domain":                validateRequired,
	"core/data-root":             validateAbsolutePath,
	"pihole/dns-port":            validatePiholeDNSPort,
	"pihole/admin-password":      validateRequired,
	"plex/media-movies":          validateOptionalAbsPath,
	"plex/media-tv":              validateOptionalAbsPath,
	"plex/media-music":           validateOptionalAbsPath,
	"nextcloud/admin-password":   validateRequired,
	"mariadb/root-password":      validateRequired,
	"mariadb/nextcloud-password": validateRequired,
}

func validateRequired(v config.Value, _ config.TreeReader) ([]config.ValidationResult, error) {
	s, _ := v.Val.(string)
	if s == "" {
		return []config.ValidationResult{{
			Message:  "This field is required",
			Severity: config.Blocking,
		}}, nil
	}
	return nil, nil
}

func validateAbsolutePath(v config.Value, _ config.TreeReader) ([]config.ValidationResult, error) {
	s, _ := v.Val.(string)
	if !strings.HasPrefix(s, "/") {
		return []config.ValidationResult{{
			Message:  "Must be an absolute path",
			Severity: config.Blocking,
		}}, nil
	}
	return nil, nil
}

func validateOptionalAbsPath(v config.Value, _ config.TreeReader) ([]config.ValidationResult, error) {
	s, _ := v.Val.(string)
	if s != "" && !strings.HasPrefix(s, "/") {
		return []config.ValidationResult{{
			Message:  "Path must be absolute",
			Severity: config.Blocking,
		}}, nil
	}
	return nil, nil
}

func validatePiholeDNSPort(v config.Value, _ config.TreeReader) ([]config.ValidationResult, error) {
	s := fmt.Sprintf("%v", v.Val)
	port, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return []config.ValidationResult{{
			Message:  "DNS port must be a number",
			Severity: config.Blocking,
		}}, nil
	}
	if port == 53 {
		return []config.ValidationResult{{
			Message:  "Port 53 may conflict with systemd-resolved. Run: sudo sed -i 's/#DNSStubListener=yes/DNSStubListener=no/' /etc/systemd/resolved.conf && sudo systemctl restart systemd-resolved",
			Severity: config.Warning,
		}}, nil
	}
	return nil, nil
}
