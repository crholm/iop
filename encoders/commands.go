package encoders

import (
	"github.com/urfave/cli/v3"
)

var Commands = []*cli.Command{
	{
		Name:   "url",
		Usage:  "url query encodes a data",
		Action: urlEncode,
	},
	{
		Name:    "binary",
		Aliases: []string{"bin", "0b"},
		Usage:   "encodes data into a binary string",
		Action:  binaryEncode,
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
		Action: base64Encode,
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
		Action: base32Encode,
	},
	{
		Name:    "hex",
		Aliases: []string{"0x"},
		Usage:   "hex encodes a data",
		Action:  hexEncode,
	},

	{
		Name:  "mime",
		Usage: "encode text as mime headers RFC 2047, ascii representations of encoded words",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "charset",
				Aliases: []string{"c"},
				Value:   "utf-8",
			},
			&cli.StringFlag{
				Name:  "schema",
				Value: "b",
				Usage: "schema to use, 'b' for base64 or 'q' for quoted printable",
			},
		},
		Action: encodeMIME,
	},
}
