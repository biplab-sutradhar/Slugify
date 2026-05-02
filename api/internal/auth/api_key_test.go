package auth

import (
	"strings"
	"testing"
)

func TestGenerateAPIKey_LengthAndAlphabet(t *testing.T) {
	key, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey error: %v", err)
	}
	if len(key) != 43 {
		t.Errorf("key length = %d, want 43", len(key))
	}
	// base64.RawURLEncoding uses A–Z a–z 0–9 - _
	allowed := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	for _, ch := range key {
		if !strings.ContainsRune(allowed, ch) {
			t.Errorf("unexpected char %q in key %q", ch, key)
		}
	}
}

func TestGenerateAPIKey_Unique(t *testing.T) {
	seen := make(map[string]struct{})
	for i := 0; i < 50; i++ {
		k, err := GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey error: %v", err)
		}
		if _, dup := seen[k]; dup {
			t.Fatalf("duplicate key generated: %q", k)
		}
		seen[k] = struct{}{}
	}
}

func TestValidateAPIKey(t *testing.T) {
	good, _ := GenerateAPIKey()

	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid generated key", good, false},
		{"empty", "", true},
		{"too short", "abc", true},
		{"too long", strings.Repeat("a", 100), true},
		{"off by one short", strings.Repeat("a", 42), true},
		{"off by one long", strings.Repeat("a", 44), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateAPIKey(c.input)
			if (err != nil) != c.wantErr {
				t.Errorf("ValidateAPIKey(%q) err = %v, wantErr = %v", c.input, err, c.wantErr)
			}
		})
	}
}
