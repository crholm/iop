package encoders

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/urfave/cli/v3"
	"io"
	"net/url"
)

func urlEncode(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	s := url.QueryEscape(string(b))
	_, err = out.Write([]byte(s))

	return err
}

func binaryEncode(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	bs, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	for _, b := range bs {
		_, err = out.Write([]byte(fmt.Sprintf("%08b", b)))
		if err != nil {
			return err
		}
	}
	return err
}

func base64Encode(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	e := base64.StdEncoding
	if c.Bool("url") {
		e = base64.URLEncoding
	}
	encoder := base64.NewEncoder(e, out)
	_, err := io.Copy(encoder, in)
	if err != nil {
		return err
	}
	return encoder.Close()
}

func base32Encode(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	e := base32.StdEncoding
	if c.Bool("hex") {
		e = base32.HexEncoding
	}
	encoder := base32.NewEncoder(e, out)
	_, err := io.Copy(encoder, in)
	if err != nil {
		return err
	}
	return encoder.Close()
}

func hexEncode(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	d := hex.NewEncoder(out)
	_, err := io.Copy(d, in)
	return err
}
