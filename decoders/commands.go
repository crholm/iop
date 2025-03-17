package decoders

import (
	"github.com/urfave/cli/v3"
)

var Commands = []*cli.Command{
	{
		Name:   "url",
		Usage:  "decodes a url query string",
		Action: decodeURL,
	},
	{
		Name:    "binary",
		Aliases: []string{"bin", "0b"},
		Usage:   "decodes a string of 1s and 0s into binary data",
		Action:  decodeBinary,
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
		Action: decodeBase64,
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
		Action: decodeBase32,
	},
	{
		Name:    "hex",
		Aliases: []string{"0x"},
		Usage:   "decodes a hex string",
		Action:  decodeHex,
	},
	{
		Name:   "jwt",
		Usage:  "decodes a jwt token",
		Action: decodeJWT,
	},
	{
		Name: "xid",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Value: "json",
				Usage: "output format, [json | yaml | xml ]",
			},
		},

		Usage:  "decodes a jwt token",
		Action: decodeXID,
	},
	{
		Name:   "mime",
		Usage:  "decodes mime headers RFC 2047, ascii representations of encoded words",
		Action: decodeMIME,
	},
}
