package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-xmlfmt/xmlfmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/satori/go.uuid"
	"github.com/urfave/cli/v2" // imports as package "cli"
	"io"
	"io/ioutil"
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
						Name:  "b64",
						Usage: "base64",
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
							if c.Bool("url"){
								e = base64.URLEncoding
							}
							d := base64.NewDecoder(e, in)
							 _, err := io.Copy(out, d)
							return err
						},
					},
					{
						Name:  "b32",
						Usage: "base32",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "hex",
								Usage:   "hex version",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							e := base32.StdEncoding
							if c.Bool("hex"){
								e = base32.HexEncoding
							}
							d := base32.NewDecoder(e, in)
							_, err := io.Copy(out, d)
							return err
						},
					},
					{
						Name:  "hex",
						Usage: "hex",
						Action: func(c *cli.Context) error {
							in := os.Stdin
							out := os.Stdout
							d := hex.NewDecoder(in)
							_, err := io.Copy(out, d)
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
						Name:  "b64",
						Usage: "base64",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "url",
								Usage:   "url encoding",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							e := base64.StdEncoding
							if c.Bool("url"){
								e = base64.URLEncoding
							}
							encoder := base64.NewEncoder(e, out)
							_, err := io.Copy(encoder, in)
							if err != nil{
								return err
							}
							return encoder.Close()
						},
					},
					{
						Name:  "b32",
						Usage: "base32",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "hex",
								Usage:   "hex version",
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							e := base32.StdEncoding
							if c.Bool("hex"){
								e = base32.HexEncoding
							}
							encoder := base32.NewEncoder(e, out)
							_, err := io.Copy(encoder, in)
							if err != nil{
								return err
							}
							return encoder.Close()
						},
					},
					{
						Name:  "hex",
						Usage: "hex",
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
						Name:  "json",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "indent",
								Value: 2,
							},
							&cli.BoolFlag{
								Name:    "color",
							},

						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							indent := c.Int("indent")
							color := c.Bool("color")

							b, err := ioutil.ReadAll(in)
							if err != nil{
								return err
							}

							f := prettyjson.NewFormatter()
							f.Indent = indent
							f.DisabledColor = !color
							b, err = f.Format(b)
							if err != nil{
								return err
							}
							_, err = io.Copy(out, bytes.NewBuffer(b))
							return err
						},
					},
					{
						Name:  "xml",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "indent",
								Value: 2,
							},
						},
						Action: func(c *cli.Context) error {

							in := os.Stdin
							out := os.Stdout

							indent := c.Int("indent")

							b, err := ioutil.ReadAll(in)

							if err != nil{
								return err
							}

							ii := strings.Join(make([]string, indent+1)," ")

							buf := xmlfmt.FormatXML(string(b), "", ii)

							_, err = io.Copy(out, bytes.NewBuffer([]byte(buf)))

							return err
						},
					},
				},
			},
			{
				Name:  "rand",
				Usage: "generate something",
				Subcommands: []*cli.Command{
					{
						Name:  "uuid",
						Usage: "generate a random v4 uuid",
						Action: func(c *cli.Context) error {
							u := uuid.NewV4()
							fmt.Println(u.String())
							return nil
						},
					},
					{
						Name:  "string",
						Usage: "generate a random string",
						Action: func(c *cli.Context) error {
							l := 20

							ls := c.Args().Get(0)
							ii, err := strconv.ParseInt(ls, 10, 32)
							if err == nil{
								l = int(ii)
							}

							var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

							// creating an even cyclic group containing all letters the same number of times
							max := len(runes) * (256/len(runes))
							acc := ""
							for {
								if len(acc) == l{
									break
								}
								b := make([]byte, 1,1)
								_, err := rand.Read(b)
								if err != nil{
									return err
								}
								// ignoring random things that is not in cyclic group,
								// since this will cause some letters to be more frequent then others
								if b[0] < byte(max){
									acc += string(runes[int(b[0]) % len(runes)])
								}
							}
							fmt.Println(acc)
							return nil
						},
					},
					{
						Name:  "bytes",
						Usage: "generate random bytes",
						Action: func(c *cli.Context) error {
							l := 20
							ls := c.Args().Get(0)
							ii, err := strconv.ParseInt(ls, 10, 32)
							if err == nil{
								l = int(ii)
							}

							var dest = make([]byte, l, l)

							n, err := rand.Read(dest)
							if err != nil{
								return err
							}
							if n != len(dest){
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
		_, _ =fmt.Fprintln(os.Stderr, "got err", err)
		os.Exit(1)
	}
}



