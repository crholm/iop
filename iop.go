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
	"math/big"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)


var out io.Writer
func main() {


	app := &cli.App{
		Name:  "iop",
		Usage: "a tool for converting and formatting things from std in to std out",
		UsageText: "You can use -- as piping between commands, eg. echo 124 | iop conv string-to-int -- encode hex -- clip copy",
		Commands: []*cli.Command{
			{
				Name:    "decode",
				Aliases: []string{"dec"},
				Usage:   "decode std from something",
				Subcommands: []*cli.Command{
					{
						Name:  "url",
						Usage: "decodes a url query string",
						Action: func(c *cli.Context) error {

							in := os.Stdin

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
						Name:    "binary",
						Aliases: []string{"bin", "0b"},
						Usage:   "decodes a string of 1s and 0s into binary data",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							bs, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}

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
						Name:    "hex",
						Aliases: []string{"0x"},
						Usage:   "decodes a hex string",
						Action: func(c *cli.Context) error {
							in := os.Stdin
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
							_, err = out.Write(tokenData)

							return err
						},
					},
				},
			},
			{
				Name:    "encode",
				Aliases: []string{"enc"},
				Usage:   "encode std to something",
				Subcommands: []*cli.Command{
					{
						Name:  "url",
						Usage: "url query encodes a data",
						Action: func(c *cli.Context) error {

							in := os.Stdin

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
						Name:    "binary",
						Aliases: []string{"bin", "0b"},
						Usage:   "encodes data into a binary string",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							bs, err := ioutil.ReadAll(in)
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
						Name:    "hex",
						Aliases: []string{"0x"},
						Usage:   "hex encodes a data",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							d := hex.NewEncoder(out)
							_, err := io.Copy(d, in)
							return err
						},
					},
				},
			},
			{
				Name:    "fmt",
				Aliases: []string{"format"},
				Usage:   "format something from std in",
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

							_, err = out.Write([]byte(str))
							return err
						},
					},
				},
			},

			{
				Name:    "rand",
				Aliases: []string{"random"},
				Usage:   "generate something random",
				Subcommands: []*cli.Command{
					{
						Name:  "uuid",
						Usage: "generate a random v4 uuid",
						Action: func(c *cli.Context) error {
							_, err := out.Write([]byte(uuid.NewV4().String()))
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
							_, err = out.Write([]byte(acc))
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
							_, err = out.Write(dest)

							return err
						},
					},
				},
			},


			{
				Name:    "conv",
				Aliases: []string{"convert"},
				Usage:   "convert something",
				Subcommands: []*cli.Command{
					{
						Name:    "int-to-string",
						Aliases: []string{"i2s"},
						Usage:   "converts byte to a string representing the number",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							bs, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}
							i := big.NewInt(0)
							i.SetBytes(bs)
							_, err = out.Write([]byte(i.String()))

							return err
						},
					},
					{
						Name:    "string-to-int",
						Aliases: []string{"s2i"},
						Usage:   "converts a string containing a number to bytes representing the number",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							bs, err := ioutil.ReadAll(in)
							if err != nil {
								return err
							}

							f, _, err := big.ParseFloat(string(bytes.TrimSpace(bs)), 10, 0, big.ToZero)
							if err != nil {
								return err
							}
							i, _ := f.Int(nil)
							_, err = out.Write(i.Bytes())

							return err
						},
					},
				},
			},
		},
	}


	var args = os.Args
	var leftover []string
	for i, a := range os.Args {
		if a == "--" {
			args = os.Args[:i]
			leftover = os.Args[i+1:]
			break
		}
	}

	out = os.Stdout

	var wg sync.WaitGroup
	if len(leftover) > 0 {
		wg.Add(1)

		cmd := exec.Command("/proc/self/exe", leftover...)
		r, w, err := os.Pipe()
		if err != nil{
			panic(err)
		}
		out = w
		cmd.Stdin = r
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func() {
			defer wg.Done()
			err := cmd.Start()

			if err != nil{
				os.Stderr.Write([]byte(err.Error()))
			}

		}()
	}
	err := app.Run(args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "got err", err)
		os.Exit(1)
	}
	wg.Wait()
}
