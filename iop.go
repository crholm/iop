package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	app := &cli.App{
		Name:  "iop",
		Usage: "a tool for converting and formatting things from std in to std out",
		Commands: []*cli.Command{
			{
				Name:  "decode",
				Usage: "decode std from something",
				Subcommands: []*cli.Command{
					{
						Name:  "url",
						Usage: "decodes a url query string",
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							b, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}
							s, err := url.QueryUnescape(string(b))
							if err != nil {
								return err
							}
							_, err = out.Write([]byte(s))

							return err
						},
					},
					{
						Name:    "b64",
						Aliases: []string{"base64"},
						Usage:   "decodes a base64 string",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "url",
								Aliases: []string{"u"},
								Usage:   "url encoding",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							e := base64.StdEncoding
							if c.Bool("url") {
								e = base64.URLEncoding
							}
							d := base64.NewDecoder(e, in)
							_, err := io.Copy(out, d)
							return err
						},
					},
					{
						Name:    "b32",
						Aliases: []string{"base32"},
						Usage:   "decodes a base32 string",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "hex",
								Usage: "hex version",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							e := base32.StdEncoding
							if c.Bool("hex") {
								e = base32.HexEncoding
							}
							d := base32.NewDecoder(e, in)
							_, err := io.Copy(out, d)
							return err
						},
					},
					{
						Name:  "hex",
						Usage: "decodes a hex string",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							out := os.Stdout
							d := hex.NewDecoder(in)
							_, err := io.Copy(out, d)
							return err
						},
					},
					{
						Name:  "jwt",
						Usage: "decodes a jwt token",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							b, err := ioutil.ReadAll(in)
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
							_, err = os.Stdout.Write(tokenData)

							return err
						},
					},
				},
			},
			{
				Name:  "encode",
				Usage: "encode std to something",
				Subcommands: []*cli.Command{
					{
						Name:  "url",
						Usage: "url query encodes a data",
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							b, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}
							s := url.QueryEscape(string(b))
							_, err = out.Write([]byte(s))

							return err
						},
					},
					{
						Name:    "b64",
						Aliases: []string{"base64"},
						Usage:   "base64 encodes a data",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "url",
								Usage: "url encoding",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

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
						},
					},
					{
						Name:    "b32",
						Aliases: []string{"base32"},
						Usage:   "base32 encodes a data",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "hex",
								Usage: "hex version",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

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
						},
					},
					{
						Name:  "hex",
						Usage: "hex encodes a data",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							out := os.Stdout
							d := hex.NewEncoder(out)
							_, err := io.Copy(d, in)
							return err
						},
					},
				},
			},
			{
				Name:  "fmt",
				Usage: "format something from std in",
				Subcommands: []*cli.Command{
					{
						Name: "json",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:  "indent",
								Value: 2,
							},
							&cli.BoolFlag{
								Name: "color",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							indent := c.Int("indent")
							color := c.Bool("color")

							b, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}

							f := prettyjson.NewFormatter()
							f.Indent = indent
							f.DisabledColor = !color
							b, err = f.Format(b)
							if err != nil {
								return err
							}
							_, err = io.Copy(out, bytes.NewBuffer(b))
							return err
						},
					},
					{
						Name: "xml",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:  "indent",
								Value: 2,
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							indent := c.Int("indent")

							b, err := ioutil.ReadAll(in)

							if err != nil {
								return err
							}

							ii := strings.Join(make([]string, indent+1), " ")

							buf := xmlfmt.FormatXML(string(b), "", ii)

							_, err = io.Copy(out, bytes.NewBuffer([]byte(buf)))

							return err
						},
					},
				},
			},
			{
				Name:  "clip",
				Usage: "managing clipboard",
				Subcommands: []*cli.Command{
					{
						Name:    "copy",
						Aliases: []string{"to", "c"},
						Usage:   "puts things from std in onto the clipboard",
						Action: func(c *cli.Context) error {
							in := os.Stdin

							b, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}
							return clipboard.WriteAll(string(b))
						},
					},
					{
						Name:    "paste",
						Aliases: []string{"from", "v"},
						Usage:   "puts things in clipboard onto std out",
						Action: func(c *cli.Context) error {
							str, err := clipboard.ReadAll()
							if err != nil {
								return err
							}

							_, err = os.Stdout.Write([]byte(str))
							return err
						},
					},
				},
			},

			{
				Name:  "rand",
				Usage: "generate something random",
				Subcommands: []*cli.Command{
					{
						Name:  "uuid",
						Usage: "generate a random v4 uuid",
						Action: func(c *cli.Context) error {
							_, err := os.Stdout.Write([]byte(uuid.NewV4().String()))
							return err
						},
					},
					{
						Name:  "string",
						Usage: "generate a random string",
						Action: func(c *cli.Context) error {
							l := 20

							ls := c.Args().Get(0)
							ii, err := strconv.ParseInt(ls, 10, 32)
							if err == nil {
								l = int(ii)
							}

							var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

							// creating an even cyclic group containing all letters the same number of times
							max := len(runes) * (256 / len(runes))
							acc := ""
							for {
								if len(acc) == l {
									break
								}
								b := make([]byte, 1, 1)
								_, err := rand.Read(b)
								if err != nil {
									return err
								}
								// ignoring random things that is not in cyclic group,
								// since this will cause some letters to be more frequent then others
								if b[0] < byte(max) {
									acc += string(runes[int(b[0])%len(runes)])
								}
							}
							_, err = os.Stdout.Write([]byte(acc))
							return err
						},
					},
					{
						Name:  "bytes",
						Usage: "generate random bytes",
						Action: func(c *cli.Context) error {
							l := 20
							ls := c.Args().Get(0)
							ii, err := strconv.ParseInt(ls, 10, 32)
							if err == nil {
								l = int(ii)
							}

							var dest = make([]byte, l, l)

							n, err := rand.Read(dest)
							if err != nil {
								return err
							}
							if n != len(dest) {
								return errors.New("could not generate enough random")
							}
							_, err = os.Stdout.Write(dest)

							return err
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "got err", err)
		os.Exit(1)
	}
}
