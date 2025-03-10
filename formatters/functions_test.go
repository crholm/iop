package formatters

import (
	"bytes"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"strings"
	"testing"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		indent   int
		color    bool
		wantErr  bool
		contains string // Check if output contains this string
	}{
		{
			name:     "valid json",
			input:    `{"name":"John","age":30}`,
			indent:   2,
			color:    false,
			wantErr:  false,
			contains: "name",
		},
		{
			name:     "invalid json",
			input:    `{"name":"John","age":30`,
			indent:   2,
			color:    false,
			wantErr:  true,
			contains: "",
		},
		{
			name:     "empty json",
			input:    `{}`,
			indent:   2,
			color:    false,
			wantErr:  false,
			contains: "{}",
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
					&cli.IntFlag{Name: "indent"},
					&cli.BoolFlag{Name: "color"},
				},
				Action: formatJSON,
			}

			args := []string{""}
			args = append(args, fmt.Sprintf("--indent=%d", tt.indent))
			if tt.color {
				args = append(args, "--color")
			}

			err := cmd.Run(context.Background(), args)

			if (err != nil) != tt.wantErr {
				t.Errorf("formatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !strings.Contains(out.String(), tt.contains) {
				t.Errorf("formatJSON() output doesn't contain expected string '%s'", tt.contains)
			}
		})
	}
}

func TestFormatXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		indent   int
		wantErr  bool
		contains string
	}{
		{
			name:     "valid xml",
			input:    `<root><item>value</item></root>`,
			indent:   2,
			wantErr:  false,
			contains: "<item>",
		},
		{
			name:     "empty xml",
			input:    `<root></root>`,
			indent:   2,
			wantErr:  false,
			contains: "<root>",
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
					&cli.IntFlag{Name: "indent"},
				},
				Action: formatXML,
			}

			args := []string{""}
			args = append(args, fmt.Sprintf("--indent=%d", tt.indent))
			err := cmd.Run(context.Background(), args)

			if (err != nil) != tt.wantErr {
				t.Errorf("formatXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !strings.Contains(out.String(), tt.contains) {
				t.Errorf("formatXML() output doesn't contain expected string '%s'", tt.contains)
			}
		})
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "mixed case",
			input:    "Hello World",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "uppercase",
			input:    "HELLO WORLD",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
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

			err := toLowerCase(context.Background(), cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("toLowerCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("toLowerCase() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}

func TestToUpperCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "mixed case",
			input:    "Hello World",
			expected: "HELLO WORLD",
			wantErr:  false,
		},
		{
			name:     "already uppercase",
			input:    "HELLO WORLD",
			expected: "HELLO WORLD",
			wantErr:  false,
		},
		{
			name:     "lowercase",
			input:    "hello world",
			expected: "HELLO WORLD",
			wantErr:  false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
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

			err := toUpperCase(context.Background(), cmd)

			if (err != nil) != tt.wantErr {
				t.Errorf("toUpperCase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && out.String() != tt.expected {
				t.Errorf("toUpperCase() got = %v, want %v", out.String(), tt.expected)
			}
		})
	}
}
