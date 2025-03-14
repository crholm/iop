package decoders

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/crholm/iop/utils"
	"github.com/hokaccha/go-prettyjson"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/xid"
	"github.com/urfave/cli/v3"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v3"
	"io"
	"math/big"
	"mime"
	"net/url"
	"strings"
	"time"
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

func decodeMIME(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	d := &mime.WordDecoder{}

	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		charset = strings.ToLower(charset)
		if m, ok := utils.CharsetEncodings[charset]; ok {
			rr := transform.NewReader(input, m.NewDecoder())
			return rr, nil
		}

		charset = utils.CharsetAliases[charset]
		if m, ok := utils.CharsetEncodings[charset]; ok {
			rr := transform.NewReader(input, m.NewDecoder())
			return rr, nil
		}

		return input, nil
	}

	header, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}
	res, err := d.DecodeHeader(strings.TrimSpace(string(header)))
	if err != nil {
		return fmt.Errorf("failed to decode header: %s", err)
	}

	_, err = out.Write([]byte(res))
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

func decodeXID(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	b, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("failed to read input: %s", err)
	}

	id, err := xid.FromString(strings.TrimSpace(string(b)))
	if err != nil {
		return fmt.Errorf("failed to parse xid: %s", err)
	}

	type xid struct {
		Time    time.Time `yaml:"time" toml:"time" json:"time" xml:"time"`
		Machine int64     `yaml:"machine" toml:"machine" json:"machine" xml:"machine"`
		Pid     uint16    `yaml:"pid" toml:"pid" json:"pid" xml:"pid"`
		Counter int32     `yaml:"counter" toml:"counter" json:"counter" xml:"counter"`
	}

	ii := big.Int{}
	ii.SetBytes(id.Machine())

	uid := xid{
		Time:    id.Time().In(time.UTC),
		Machine: ii.Int64(),
		Pid:     id.Pid(),
		Counter: id.Counter(),
	}

	var marshaller func(any) ([]byte, error)

	switch c.String("format") {
	case "json":
		marshaller = json.Marshal
	case "xml":
		marshaller = xml.Marshal
	case "yaml", "yml":
		marshaller = yaml.Marshal
	case "toml":
		marshaller = toml.Marshal
	case "text", "txt":
		marshaller = func(z any) ([]byte, error) {
			return []byte(fmt.Sprintf("%v", z)), nil
		}
	default:
		marshaller = json.Marshal
	}

	j, err := marshaller(uid)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}
	_, err = out.Write(j)

	return err

}
