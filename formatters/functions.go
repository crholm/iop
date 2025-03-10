package formatters

import (
	"bytes"
	"context"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli/v3"
	"io"
	"strings"
)

func formatJSON(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	indent := c.Int("indent")
	color := c.Bool("color")

	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	f := prettyjson.NewFormatter()
	f.Indent = int(indent)
	f.DisabledColor = !color
	b, err = f.Format(b)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, bytes.NewBuffer(b))
	return err
}

func formatXML(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	indent := c.Int("indent")

	b, err := io.ReadAll(in)

	if err != nil {
		return err
	}

	ii := strings.Join(make([]string, indent+1), " ")

	buf := xmlfmt.FormatXML(string(b), "", ii)

	_, err = io.Copy(out, bytes.NewBuffer([]byte(buf)))

	return err
}

func toLowerCase(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	b, err := io.ReadAll(in)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewBuffer(bytes.ToLower(b)))

	return err
}

func toUpperCase(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	b, err := io.ReadAll(in)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewBuffer(bytes.ToUpper(b)))

	return err
}
