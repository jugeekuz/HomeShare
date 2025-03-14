package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	const testOrigin = "https://example.com"

	t.Setenv("DOMAIN_ORIGIN", testOrigin)

	cfg := LoadConfig()

	if cfg.DomainOrigin != testOrigin {
		t.Errorf("expected DomainOrigin to be %q, got %q", testOrigin, cfg.DomainOrigin)
	}
}
