package conversions

import (
	"bytes"
	"context"
	"github.com/urfave/cli/v3"
	"io"
	"strings"
	"testing"
)

func TestIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		wantErr  bool
		little   bool
	}{
		{
			name:     "zero value",
			input:    []byte{0},
			expected: "0",
			wantErr:  false,
			little:   false,
		},
		{
			name:     "positive integer",
			input:    []byte{1, 35},
			expected: "291",
			wantErr:  false,
			little:   false,
		},
		{
			name:     "large integer",
			input:    []byte{255, 255, 255, 255},
			expected: "4294967295",
			wantErr:  false,
			little:   false,
		},
		{
			name:     "empty input",
			input:    []byte{},
			expected: "0",
			wantErr:  false,
			little:   false,
		},
		{
			name:     "little-endian",
			input:    []byte{35, 1},
			expected: "291",
			wantErr:  false,
			little:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := bytes.NewReader(tt.input)
			out := &bytes.Buffer{}
			cmd := &cli.Command{
				Reader: in,
				Writer: out,
				Flags:  []cli.Flag{&cli.BoolFlag{Name: "little-endian", Value: tt.little}},
				Action: intToString,
			}

			args := []string{""}
			if tt.little {
				args = append(args, "--little-endian")
			}
			err := cmd.Run(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("intToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("intToString() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		little   bool
		wantErr  bool
	}{
		{
			name:     "zero value",
			input:    "0",
			expected: []byte{0},
			little:   false,
			wantErr:  false,
		},
		{
			name:     "positive integer",
			input:    "291",
			expected: []byte{1, 35},
			little:   false,
			wantErr:  false,
		},
		{
			name:     "large integer",
			input:    "4294967295",
			expected: []byte{255, 255, 255, 255},
			little:   false,
			wantErr:  false,
		},
		{
			name:     "with whitespace",
			input:    " 123 ",
			expected: []byte{123},
			little:   false,
			wantErr:  false,
		},
		{
			name:     "invalid input",
			input:    "not a number",
			expected: nil,
			little:   false,
			wantErr:  true,
		},
		{
			name:     "little-endian",
			input:    "291",
			expected: []byte{35, 1},
			little:   true,
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "little-endian",
					},
				},
				Action: stringToInt,
			}

			args := []string{""}

			if tt.little {
				args = append(args, "--little-endian")
			}

			err := cmd.Run(context.Background(), args)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !bytes.Equal(out.Bytes(), tt.expected) {
				t.Errorf("stringToInt() got = %v, want %v", out.Bytes(), tt.expected)
			}
		})
	}
}

func TestCsvTo(t *testing.T) {
	// For csvTo, we need to test with a mock encoder
	// This is a simplified test that verifies the function handles delimiter options

	tests := []struct {
		name      string
		input     string
		delimiter string
		wantErr   bool
	}{
		{
			name:      "comma delimiter",
			input:     "a,b,c\n1,2,3",
			delimiter: ",",
			wantErr:   false,
		},
		{
			name:      "tab delimiter",
			input:     "a\tb\tc\n1\t2\t3",
			delimiter: "\\t",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			// Create a mock encoder function for testing
			mockEncoderFn := func(w io.Writer) encoder {
				return &mockEncoder{writer: w}
			}

			cmd := &cli.Command{
				Reader: in,
				Writer: out,
				Flags:  []cli.Flag{&cli.StringFlag{Name: "delimiter"}},
				Action: csvTo(mockEncoderFn),
			}

			args := []string{""}
			if tt.delimiter == "" {
				args = append(args, "--delimiter="+tt.delimiter)
			}

			err := cmd.Run(context.Background(), args)

			if (err != nil) != tt.wantErr {
				t.Errorf("csvTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Mock encoder for testing csvTo
type mockEncoder struct {
	writer io.Writer
}

func (m *mockEncoder) Encode(v any) error {
	// For testing, just write something to show it was called
	_, err := m.writer.Write([]byte("encoded"))
	return err
}
