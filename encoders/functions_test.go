package encoders

import (
	"bytes"
	"context"
	"github.com/urfave/cli/v3"
	"strings"
	"testing"
)

func TestURLEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple url encode",
			input:    "hello world",
			expected: "hello+world",
			wantErr:  false,
		},
		{
			name:     "encode special characters",
			input:    "param1=value1&param2=value2",
			expected: "param1%3Dvalue1%26param2%3Dvalue2",
			wantErr:  false,
		},
		{
			name:     "encode symbols",
			input:    "!@#$%^&*()",
			expected: "%21%40%23%24%25%5E%26%2A%28%29",
			wantErr:  false,
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
			err := urlEncode(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("urlEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("urlEncode() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestBinaryEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple binary encode",
			input:    "A",
			expected: "01000001",
			wantErr:  false,
		},
		{
			name:     "encode Hello",
			input:    "Hello",
			expected: "0100100001100101011011000110110001101111",
			wantErr:  false,
		},
		{
			name:     "encode with special chars",
			input:    "A!",
			expected: "0100000100100001",
			wantErr:  false,
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
			err := binaryEncode(context.Background(), cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("binaryEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("binaryEncode() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		urlMode  bool
		wantErr  bool
	}{
		{
			name:     "standard base64 encode",
			input:    "Hello World",
			expected: "SGVsbG8gV29ybGQ=",
			urlMode:  false,
			wantErr:  false,
		},
		{
			name:     "url base64 encode",
			input:    "Hello World?",
			expected: "SGVsbG8gV29ybGQ_",
			urlMode:  true,
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			urlMode:  false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "url", Value: tt.urlMode}},
				Action: base64Encode,
			}
			var args []string = []string{""}
			// Set the url flag if needed
			if tt.urlMode {
				args = append(args, "--url=true")
			}
			err := cmd.Run(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("base64Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("base64Encode() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestBase32Encode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hexMode  bool
		wantErr  bool
	}{
		{
			name:     "standard base32 encode",
			input:    "Hello World",
			expected: "JBSWY3DPEBLW64TMMQ======",
			hexMode:  false,
			wantErr:  false,
		},
		{
			name:     "hex base32 encode",
			input:    "Hello World",
			expected: "91IMOR3F41BMUSJCCG======",
			hexMode:  true,
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			hexMode:  false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "hex"}},
				Action: base32Encode,
			}

			args := []string{""}
			// Set the hex flag if needed
			if tt.hexMode {
				args = append(args, "--hex=true")
			}

			err := cmd.Run(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("base32Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("base32Encode() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestHexEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple hex encode",
			input:    "Hello World",
			expected: "48656c6c6f20576f726c64",
			wantErr:  false,
		},
		{
			name:     "encode special characters",
			input:    "!@#$",
			expected: "21402324",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "binary data",
			input:    string([]byte{0x00, 0xFF, 0x7F}),
			expected: "00ff7f",
			wantErr:  false,
		},
		{
			name:     "unicode characters",
			input:    "こんにちは",
			expected: "e38193e38293e381abe381a1e381af",
			wantErr:  false,
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

			err := hexEncode(context.Background(), cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("hexEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("hexEncode() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}
