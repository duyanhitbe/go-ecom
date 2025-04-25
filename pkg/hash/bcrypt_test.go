package hash

import (
	"testing"
)

func TestBcrypt_Hash(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		expectErr bool
	}{
		{"valid password", "securepassword", false},
		{"empty password", "", false},
		{"long password", string(make([]byte, 1000)), true},
	}

	bcrypt := NewBcrypt()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := bcrypt.Hash(tt.password)
			if (err != nil) != tt.expectErr {
				t.Errorf("Hash() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if len(hash) == 0 && !tt.expectErr {
				t.Errorf("Hash() returned empty hash for password: %v", tt.password)
			}
		})
	}
}

func TestBcrypt_Verify(t *testing.T) {
	bcrypt := NewBcrypt()
	hash, _ := bcrypt.Hash("securepassword")

	tests := []struct {
		name      string
		hash      string
		password  string
		expectOk  bool
		expectErr bool
	}{
		{"valid hash and password", hash, "securepassword", true, false},
		{"invalid password", hash, "wrongpassword", false, true},
		{"empty hash", "", "securepassword", false, true},
		{"empty password", hash, "", false, true},
		{"malformed hash", "invalidhash", "securepassword", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := bcrypt.Verify(tt.hash, tt.password)
			if ok != tt.expectOk {
				t.Errorf("Verify() ok = %v, expectOk %v", ok, tt.expectOk)
			}
			if (err != nil) != tt.expectErr {
				t.Errorf("Verify() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}
