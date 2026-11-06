// Package config loads and validates deployment configuration.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"unicode"

	"github.com/tarangchikhalia/authopus/internal/identity"
)

type Identity struct {
	Version           int       `json:"version"`
	OrganizationID    string    `json:"organization_id"`
	Issuer            string    `json:"issuer"`
	OrganizationClaim string    `json:"organization_claim,omitempty"`
	Audiences         Audiences `json:"audiences"`
	AllowedAlgorithms []string  `json:"allowed_algorithms"`
	ClockSkewSeconds  int       `json:"clock_skew_seconds"`
	Bindings          []Binding `json:"bindings"`
}
type Audiences struct {
	Caller []string `json:"caller"`
	Agent  []string `json:"agent"`
}
type Binding struct {
	Subject          string        `json:"subject"`
	ClientID         string        `json:"client_id"`
	Kind             identity.Kind `json:"kind"`
	AssertableAgents []string      `json:"assertable_agents,omitempty"`
}
type ValidationError struct {
	Path string
	Code string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("identity configuration: %s: %s", e.Path, e.Code)
}
func invalid(path, code string) error { return &ValidationError{path, code} }

// checkJSON rejects duplicate keys and nulls before typed decoding.
func checkJSON(d *json.Decoder) error {
	t, err := d.Token()
	if err != nil {
		return invalid("$", "invalid_json")
	}
	if t == nil {
		return invalid("$", "null_value")
	}
	if v, ok := t.(json.Delim); ok {
		switch v {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				if e != nil {
					return invalid("$", "invalid_json")
				}
				key, ok := k.(string)
				if !ok || seen[key] {
					return invalid("$", "duplicate_key")
				}
				seen[key] = true
				if e = checkJSON(d); e != nil {
					return e
				}
			}
		case '[':
			for d.More() {
				if e := checkJSON(d); e != nil {
					return e
				}
			}
		default:
			return invalid("$", "invalid_json")
		}
		if _, e := d.Token(); e != nil {
			return invalid("$", "invalid_json")
		}
	}
	return nil
}
func LoadIdentity(r io.Reader) (Identity, error) {
	var c Identity
	b, e := io.ReadAll(io.LimitReader(r, (1<<20)+1))
	if e != nil {
		return c, invalid("$", "read_failed")
	}
	if len(b) > 1<<20 {
		return c, invalid("$", "too_large")
	}
	d := json.NewDecoder(strings.NewReader(string(b)))
	if e = checkJSON(d); e != nil {
		return c, e
	}
	if _, e = d.Token(); e != io.EOF {
		return c, invalid("$", "trailing_json")
	}
	c.ClockSkewSeconds = 30
	d = json.NewDecoder(strings.NewReader(string(b)))
	d.DisallowUnknownFields()
	if e = d.Decode(&c); e != nil {
		return Identity{}, invalid("$", "invalid_shape")
	}
	if e = c.Validate(); e != nil {
		return Identity{}, e
	}
	return c, nil
}
func nonempty(s string) bool { return strings.TrimSpace(s) != "" }
func unique(values []string) bool {
	seen := map[string]bool{}
	for _, v := range values {
		if !nonempty(v) || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}
func (c Identity) Validate() error {
	if c.Version != 1 {
		return invalid("version", "unsupported_version")
	}
	if !nonempty(c.OrganizationID) {
		return invalid("organization_id", "required")
	}
	u, e := url.Parse(c.Issuer)
	if e != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.Contains(c.Issuer, "#") {
		return invalid("issuer", "invalid_https_issuer")
	}
	if c.OrganizationClaim != "" {
		for _, r := range c.OrganizationClaim {
			if unicode.IsSpace(r) || unicode.IsControl(r) {
				return invalid("organization_claim", "invalid_claim_name")
			}
		}
	}
	for _, a := range []struct {
		p string
		v []string
	}{{"audiences.caller", c.Audiences.Caller}, {"audiences.agent", c.Audiences.Agent}, {"allowed_algorithms", c.AllowedAlgorithms}} {
		if len(a.v) == 0 || !unique(a.v) {
			return invalid(a.p, "required_unique_values")
		}
	}
	for _, a := range c.AllowedAlgorithms {
		if a != "RS256" && a != "ES256" {
			return invalid("allowed_algorithms", "unsupported_algorithm")
		}
	}
	if c.ClockSkewSeconds < 0 || c.ClockSkewSeconds > 60 {
		return invalid("clock_skew_seconds", "out_of_range")
	}
	if len(c.Bindings) == 0 {
		return invalid("bindings", "required")
	}
	subjects := map[string]identity.Kind{}
	clients := map[string]bool{}
	for i, b := range c.Bindings {
		p := fmt.Sprintf("bindings[%d]", i)
		if !nonempty(b.Subject) || !nonempty(b.ClientID) {
			return invalid(p, "required_identity")
		}
		if b.Kind != identity.Agent && b.Kind != identity.Integration {
			return invalid(p+".kind", "unknown_kind")
		}
		if _, ok := subjects[b.Subject]; ok || clients[b.ClientID] {
			return invalid(p, "duplicate_binding")
		}
		subjects[b.Subject] = b.Kind
		clients[b.ClientID] = true
	}
	for i, b := range c.Bindings {
		p := fmt.Sprintf("bindings[%d].assertable_agents", i)
		if b.Kind != identity.Integration && len(b.AssertableAgents) > 0 {
			return invalid(p, "integration_only")
		}
		if !unique(b.AssertableAgents) {
			return invalid(p, "duplicate_target")
		}
		for _, s := range b.AssertableAgents {
			if subjects[s] != identity.Agent {
				return invalid(p, "unknown_agent")
			}
		}
	}
	return nil
}
