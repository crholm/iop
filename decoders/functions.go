package decoders

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli/v3"
	"io"
	"net/url"
	"strings"
)

func decodeURL(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	s, err := url.QueryUnescape(string(b))
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(s))

	return err
}

func decodeBinary(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	bs, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	bs = bytes.ReplaceAll(bs, []byte(" "), []byte(""))
	var bb byte
	for i, b := range bytes.TrimSpace(bs) {

		bb = bb << 1
		switch b {
		case '1':
			bb = bb ^ 1
		case '0':
		default:
			return errors.New("non 1,0 char was found")
		}
		if i%8 == 7 {
			_, err = out.Write([]byte{bb})
			if err != nil {
				return err
			}
			bb = 0
		}
	}
	return err
}

func decodeBase64(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	e := base64.StdEncoding
	if c.Bool("url") {
		e = base64.URLEncoding
	}
	d := base64.NewDecoder(e, in)
	_, err := io.Copy(out, d)
	return err
}

func decodeBase32(ctx context.Context, c *cli.Command) error {

	in := c.Reader
	out := c.Writer

	e := base32.StdEncoding
	if c.Bool("hex") {
		e = base32.HexEncoding
	}
	d := base32.NewDecoder(e, in)
	_, err := io.Copy(out, d)
	return err
}

func decodeHex(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	d := hex.NewDecoder(in)
	_, err := io.Copy(out, d)
	return err
}

func decodeJWT(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	b = bytes.ReplaceAll(b, []byte(" "), []byte(""))
	b = bytes.ReplaceAll(b, []byte("\n"), []byte(""))
	b = bytes.ReplaceAll(b, []byte("\r"), []byte(""))

	parts := strings.Split(string(b), ".")
	if len(parts) != 3 {
		return errors.New("expected 3 parts of the jwt token, there were " + fmt.Sprintf("%d", len(parts)))
	}
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return err
	}
	sigString := strings.TrimSpace(parts[2])
	token := struct {
		Header    json.RawMessage `json:"header"`
		Payload   json.RawMessage `json:"payload"`
		Signature string          `json:"signature"`
	}{
		Header:    header,
		Payload:   payload,
		Signature: sigString,
	}

	f := prettyjson.NewFormatter()
	f.Indent = 2
	f.DisabledColor = false
	tokenData, err := f.Marshal(token)
	if err != nil {
		return err
	}
	_, err = out.Write(tokenData)

	return err
}
