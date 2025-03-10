package generators

import (
	"bytes"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestGeneratePassphrase(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		short     bool
		mix       bool
		wordCount int
		wantErr   bool
	}{
		{
			name:      "default passphrase",
			args:      []string{""},
			short:     false,
			mix:       false,
			wordCount: 4, // Default word count
			wantErr:   false,
		},
		{
			name:      "custom word count",
			args:      []string{"6"},
			short:     false,
			mix:       false,
			wordCount: 6,
			wantErr:   false,
		},
		{
			name:      "short flag",
			args:      []string{""},
			short:     true,
			mix:       false,
			wordCount: 4,
			wantErr:   false,
		},
		{
			name:      "mix flag",
			args:      []string{""},
			short:     false,
			mix:       true,
			wordCount: 4,
			wantErr:   false,
		},
		{
			name:      "invalid word count",
			args:      []string{"abc"},
			short:     false,
			mix:       false,
			wordCount: 4, // Should fall back to default
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Name:   "passphrase",
				Writer: out,
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "short"},
					&cli.BoolFlag{Name: "mix"},
				},
				Action: generatePassphrase,
			}

			args := append([]string{"passphrase"}, tt.args...)
			if tt.short {
				args = append([]string{"--short"}, args...)
			}
			if tt.mix {
				args = append([]string{"--mix"}, args...)
			}

			err := cmd.Run(context.Background(), args)

			if (err != nil) != tt.wantErr {
				t.Errorf("generatePassphrase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := out.String()
			words := strings.Split(result, "-")

			if len(words) != tt.wordCount {
				t.Errorf("generatePassphrase() word count = %v, want %v", len(words), tt.wordCount)
			}
		})
	}
}
func TestGenerateUUID(t *testing.T) {
	tests := []struct {
		name      string
		version   int
		namespace string
		data      string
		wantErr   bool
	}{
		{
			name:      "default version (4)",
			version:   4,
			namespace: "",
			data:      "",
			wantErr:   false,
		},
		{
			name:      "version 3 with namespace and data",
			version:   3,
			namespace: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", // DNS namespace UUID
			data:      "example.com",
			wantErr:   false,
		},
		{
			name:      "version 4 random",
			version:   4,
			namespace: "",
			data:      "",
			wantErr:   false,
		},
		{
			name:      "version 5 with namespace and data",
			version:   5,
			namespace: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", // DNS namespace UUID
			data:      "example.com",
			wantErr:   false,
		},
		{
			name:      "version 6",
			version:   6,
			namespace: "",
			data:      "",
			wantErr:   false,
		},
		{
			name:      "version 7",
			version:   7,
			namespace: "",
			data:      "",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Name:   "uuid",
				Writer: out,
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "version"},
					&cli.StringFlag{Name: "namespace"},
					&cli.StringFlag{Name: "data"},
				},
				Action: generateUUID,
			}

			args := []string{cmd.Name}
			// Set flag values
			args = append(args, fmt.Sprintf("--version=%d", tt.version))

			if tt.namespace != "" {
				args = append(args, fmt.Sprintf("--namespace=%s", tt.namespace))
			}
			if tt.data != "" {
				args = append(args, fmt.Sprintf("--data=%s", tt.data))
			}

			err := cmd.Run(context.Background(), args)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				result := strings.TrimSpace(out.String())

				// Validate UUID format (8-4-4-4-12 hex digits)
				uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
				if !uuidPattern.MatchString(result) {
					t.Errorf("generateUUID() output = %v, doesn't match UUID format", result)
					return
				}

				// Check version number in the UUID
				// Version is stored in the 13th character of the UUID string (after removing dashes)
				// For UUID string format: xxxxxxxx-xxxx-Vxxx-xxxx-xxxxxxxxxxxx where V is the version
				versionChar := result[14]
				expectedVersion := '0' + byte(tt.version)

				if versionChar != expectedVersion {
					t.Errorf("generateUUID() version = %c, want %c", versionChar, expectedVersion)
				}

				// For version 3 and 5, same namespace and data should produce the same UUID
				if tt.version == 3 || tt.version == 5 {
					// Run the function again with the same inputs
					out2 := &bytes.Buffer{}
					cmd2 := &cli.Command{
						Name:   "uuid",
						Writer: out2,
						Flags:  cmd.Flags,
						Action: cmd.Action,
					}

					err = cmd2.Run(context.Background(), []string{cmd2.Name, "--version", strconv.Itoa(tt.version), "--namespace", tt.namespace, "--data", tt.data})
					if err != nil {
						t.Errorf("generateUUID() second call error = %v", err)
						return
					}

					result2 := strings.TrimSpace(out2.String())

					// The UUIDs should be identical for deterministic versions
					if result != result2 {
						t.Errorf("generateUUID() deterministic UUID mismatch: %v != %v", result, result2)
					}
				}
			}
		})
	}
}

func TestGenerateXID(t *testing.T) {
	out := &bytes.Buffer{}
	cmd := &cli.Command{
		Writer: out,
	}

	err := generateXID(context.Background(), cmd)

	if err != nil {
		t.Errorf("generateXID() error = %v", err)
		return
	}

	result := strings.TrimSpace(out.String())

	// XID is 20 characters
	if len(result) != 20 {
		t.Errorf("generateXID() output length = %v, want 20", len(result))
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		length  int
		wantErr bool
	}{
		{
			name:    "default length",
			args:    []string{},
			length:  16, // Default length
			wantErr: false,
		},
		{
			name:    "custom length",
			args:    []string{"32"},
			length:  32,
			wantErr: false,
		},
		{
			name:    "invalid length",
			args:    []string{"abc"},
			length:  16, // Should fall back to default
			wantErr: false,
		},
		{
			name:    "negative length",
			args:    []string{"\\-10"},
			length:  16, // Should fall back to default
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Name:   "random-string",
				Writer: out,
				Action: generateRandomString,
			}

			args := append([]string{"random-string"}, tt.args...)
			err := cmd.Run(context.Background(), args)

			if (err != nil) != tt.wantErr {
				t.Errorf("generateRandomString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := strings.TrimSpace(out.String())

			if len(result) != tt.length {
				t.Errorf("generateRandomString() output length = %v, want %v", len(result), tt.length)
			}
		})
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		length  int
		wantErr bool
	}{
		{
			name:    "default length",
			args:    []string{""},
			length:  16, // Default length
			wantErr: false,
		},
		{
			name:    "custom length",
			args:    []string{"32"},
			length:  32,
			wantErr: false,
		},
		{
			name:    "invalid length",
			args:    []string{"abc"},
			length:  16, // Should fall back to default
			wantErr: false,
		},
		{
			name:    "negative length",
			args:    []string{"\"-10\""},
			length:  16, // Should fall back to default
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Name:   "random-bytes",
				Writer: out,
				Action: generateRandomBytes,
			}

			err := cmd.Run(context.Background(), append([]string{"random-bytes"}, tt.args...))

			if (err != nil) != tt.wantErr {
				t.Errorf("generateRandomBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result := out.Bytes()

			// Check that we got the expected number of bytes
			if len(result) != tt.length {
				t.Errorf("generateRandomBytes() output length = %v, want %v", len(result), tt.length)
			}
		})
	}
}
