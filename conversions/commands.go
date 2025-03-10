package conversions

import (
	"encoding/json"
	"encoding/xml"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
	"io"
)

var Commands = []*cli.Command{
	{
		Name:    "csv-to-yaml",
		Aliases: []string{"csv-to-yml"},
		Usage:   "converts a csv file to yaml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "delimiter",
				Aliases: []string{"d"},
				Value:   ",",
			},
			&cli.BoolFlag{
				Name:    "with-headers",
				Aliases: []string{"H"},
			},
		},
		Action: csvTo(func(w io.Writer) encoder {
			return yaml.NewEncoder(w)
		}),
	},
	{
		Name:  "csv-to-json",
		Usage: "converts a csv file to json",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "delimiter",
				Aliases: []string{"d"},
				Value:   ",",
			},
			&cli.BoolFlag{
				Name:    "with-headers",
				Aliases: []string{"H"},
			},
		},
		Action: csvTo(func(w io.Writer) encoder {
			return json.NewEncoder(w)
		}),
	},
	{
		Name:  "csv-to-xml",
		Usage: "converts a csv file to xml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "delimiter",
				Aliases: []string{"d"},
				Value:   ",",
			},
			&cli.BoolFlag{
				Name:    "with-headers",
				Aliases: []string{"H"},
			},
		},
		Action: csvTo(func(w io.Writer) encoder {
			return xml.NewEncoder(w)
		}),
	},

	{
		Name:    "int-to-string",
		Aliases: []string{"i2s"},
		Usage:   "converts byte to a string representing the number",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "little-endian",
				Usage: "converts little endian bytes to a string. (default is big endian)",
			},
		},
		Action: intToString,
	},
	{
		Name:    "string-to-int",
		Aliases: []string{"s2i"},
		Usage:   "converts a string containing a number to bytes representing the number",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "little-endian",
				Usage: "converts a string to little endian bytes. (default is big endian)",
			},
		},
		Action: stringToInt,
	},
}
