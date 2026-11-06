package config

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestLoadIdentity(t *testing.T) {
	b, e := os.ReadFile("testdata/identity.json")
	if e != nil {
		t.Fatal(e)
	}
	base := string(b)
	cases := []struct {
		name, input string
		valid       bool
		skew        int
	}{
		{"example", base, true, 30},
		{"zero skew", strings.Replace(base, `"version": 1`, `"clock_skew_seconds": 0, "version": 1`, 1), true, 0},
		{"unknown field", strings.Replace(base, `"version"`, `"secret-sentinel"`, 1), false, 0},
		{"duplicate key", strings.Replace(base, `"version": 1`, `"version": 1, "version": 1`, 1), false, 0},
		{"null", strings.Replace(base, `"version": 1`, `"version": null`, 1), false, 0},
		{"version", strings.Replace(base, `"version": 1`, `"version": 2`, 1), false, 0},
		{"http", strings.Replace(base, "https:", "http:", 1), false, 0},
		{"algorithm", strings.Replace(base, "RS256", "HS256", 1), false, 0},
		{"kind", strings.Replace(base, `"kind": "agent"`, `"kind": "human"`, 1), false, 0},
		{"missing target", strings.Replace(base, `"assertable_agents": ["11111111-1111-4111-8111-111111111111"]`, `"assertable_agents": ["unknown"]`, 1), false, 0},
		{"duplicate client", strings.Replace(base, "integration-gateway", "agent-research", 1), false, 0},
		{"trailing", base + "{}", false, 0},
		{"too large", strings.Repeat(" ", 1<<20) + base, false, 0},
		{"claim whitespace", strings.Replace(base, `"version": 1`, `"organization_claim": "bad name", "version": 1`, 1), false, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := LoadIdentity(strings.NewReader(tc.input))
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				if c.ClockSkewSeconds != tc.skew {
					t.Fatal("incorrect default")
				}
			} else {
				var v *ValidationError
				if !errors.As(err, &v) {
					t.Fatalf("expected typed failure, got %v", err)
				}
				if strings.Contains(err.Error(), "secret-sentinel") {
					t.Fatal("input leaked")
				}
			}
		})
	}
}

func TestValidationRules(t *testing.T) {
	mutations := []struct {
		name   string
		change func(*Identity)
	}{
		{"empty organization", func(c *Identity) { c.OrganizationID = "" }},
		{"missing audiences", func(c *Identity) { c.Audiences.Agent = nil }},
		{"duplicate audience", func(c *Identity) { c.Audiences.Caller = []string{"a", "a"} }},
		{"duplicate algorithm", func(c *Identity) { c.AllowedAlgorithms = []string{"RS256", "RS256"} }},
		{"negative skew", func(c *Identity) { c.ClockSkewSeconds = -1 }},
		{"excessive skew", func(c *Identity) { c.ClockSkewSeconds = 61 }},
		{"missing bindings", func(c *Identity) { c.Bindings = nil }},
		{"duplicate subject", func(c *Identity) { c.Bindings[1].Subject = c.Bindings[0].Subject }},
		{"agent assertion", func(c *Identity) { c.Bindings[0].AssertableAgents = []string{c.Bindings[0].Subject} }},
		{"integration target", func(c *Identity) { c.Bindings[1].AssertableAgents = []string{c.Bindings[1].Subject} }},
		{"duplicate target", func(c *Identity) {
			c.Bindings[1].AssertableAgents = []string{c.Bindings[0].Subject, c.Bindings[0].Subject}
		}},
		{"issuer credentials", func(c *Identity) { c.Issuer = "https://secret@example.test/realm" }},
		{"issuer query", func(c *Identity) { c.Issuer += "?x=secret" }},
		{"issuer fragment", func(c *Identity) { c.Issuer += "#secret" }},
	}
	for _, tc := range mutations {
		t.Run(tc.name, func(t *testing.T) {
			b, err := os.ReadFile("testdata/identity.json")
			if err != nil {
				t.Fatal(err)
			}
			c, err := LoadIdentity(strings.NewReader(string(b)))
			if err != nil {
				t.Fatal(err)
			}
			tc.change(&c)
			if c.Validate() == nil {
				t.Fatal("invalid configuration accepted")
			}
		})
	}
}
