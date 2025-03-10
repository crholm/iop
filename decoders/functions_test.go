package decoders

import (
	"bytes"
	"context"
	"github.com/urfave/cli/v3"
	"strings"
	"testing"
)

func TestDecodeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple url decode",
			input:    "hello%20world",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "complex url decode",
			input:    "param1%3Dvalue1%26param2%3Dvalue2",
			expected: "param1=value1&param2=value2",
			wantErr:  false,
		},
		{
			name:     "invalid url encoding",
			input:    "invalid%",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
			}
			err := decodeURL(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("decodeURL() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		urlMode  bool
		wantErr  bool
	}{
		{
			name:     "standard base64 decode",
			input:    "SGVsbG8gV29ybGQ=",
			expected: "Hello World",
			urlMode:  false,
			wantErr:  false,
		},
		{
			name:     "url base64 decode",
			input:    "SGVsbG8gV29ybGQ=",
			expected: "Hello World",
			urlMode:  true,
			wantErr:  false,
		},
		{
			name:     "invalid base64",
			input:    "invalid-base64!@#",
			expected: "",
			urlMode:  false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "url"}},
				Action: decodeBase64,
			}

			args := []string{""}
			if tt.urlMode {
				args = append(args, "--url")
			}

			err := cmd.Run(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("decodeBase64() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestDecodeHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid hex decode",
			input:    "48656c6c6f20576f726c64",
			expected: "Hello World",
			wantErr:  false,
		},
		{
			name:     "invalid hex",
			input:    "invalid-hex!@#",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
			}
			err := decodeHex(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("decodeHex() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestDecodeBinary(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid binary decode",
			input:    "01001000 01100101 01101100 01101100 01101111", // "Hello" in binary
			expected: "Hello",
			wantErr:  false,
		},
		{
			name:     "valid binary decode no space",
			input:    "0100100001100101011011000110110001101111", // "Hello" in binary
			expected: "Hello",
			wantErr:  false,
		},
		{
			name:     "invalid binary with non-binary character",
			input:    "0102",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
			}
			err := decodeBinary(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("decodeBinary() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}
func TestDecodeJWT(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid JWT token",
			input:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantErr: false,
		},
		{
			name:    "invalid JWT format",
			input:   "invalid.jwt.token",
			wantErr: true, // Should error with invalid base64
		},
		{
			name:    "incomplete JWT",
			input:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			wantErr: true, // Missing signature part
		},
		{
			name:    "JWT with too many parts",
			input:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c.extra",
			wantErr: true, // Too many parts
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
			}

			err := decodeJWT(context.Background(), cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("decodeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify output contains expected JWT parts
				output := out.String()
				if !strings.Contains(output, "header") || !strings.Contains(output, "payload") || !strings.Contains(output, "signature") {
					t.Errorf("decodeJWT() output missing expected JWT components: %v", output)
				}
			}
		})
	}
}
