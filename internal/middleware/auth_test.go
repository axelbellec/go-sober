package middleware

import (
	"net/http"
	"testing"
)

func TestExtractTokenFromHeader(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Bearer token",
			authHeader:    "Bearer abc123.xyz789.def456",
			expectedToken: "abc123.xyz789.def456",
			expectError:   false,
		},
		{
			name:        "Missing Authorization header",
			authHeader:  "",
			expectError: true,
		},
		{
			name:        "Missing Bearer prefix",
			authHeader:  "abc123.xyz789.def456",
			expectError: true,
		},
		{
			name:        "Invalid format - extra parts",
			authHeader:  "Bearer abc123.xyz789.def456 extra",
			expectError: true,
		},
		{
			name:        "Invalid format - wrong prefix",
			authHeader:  "Basic abc123.xyz789.def456",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			token, err := extractTokenFromHeader(req)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("expected token %q, got %q", tt.expectedToken, token)
				}
			}
		})
	}
}
