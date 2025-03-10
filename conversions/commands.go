package conversions

import (
	"encoding/json"
	"encoding/xml"
	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
	"io"
)

var Commands = []*cli.Command{
	{
		Name:  "csv-to-yaml",
		Usage: "converts a csv file to yaml",
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
		Name:  "csv-to-toml",
		Usage: "converts a csv file to toml",
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
			return toml.NewEncoder(w)
		}),
	},

	// JSON-
	{
		Name:  "json-to-csv",
		Usage: `converts json to csv. It must be a list of objects, [{"a":1, "b":2}, {"a":3, "b":4}]`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "delimiter",
				Aliases: []string{"d"},
				Value:   ",",
			},
		},
		Action: toCsv(decoderJSON),
	},
	{
		Name:   "json-to-xml",
		Usage:  "converts json to xml (WARNING: works poorly, xml is broken)",
		Action: stdFromTo(decoderJSON, encoderXML),
	},

	{
		Name:   "json-to-toml",
		Usage:  "converts json to toml",
		Action: stdFromTo(decoderJSON, encoderTOML),
	},
	{
		Name:   "json-to-yaml",
		Usage:  "converts json to yaml",
		Action: stdFromTo(decoderJSON, encoderYAML),
	},

	// TOML -
	{
		Name:   "toml-to-xml",
		Usage:  "converts toml to xml (WARNING: works poorly, xml is broken)",
		Action: stdFromTo(decoderTOML, encoderXML),
	},
	{
		Name:   "toml-to-json",
		Usage:  "converts toml to json",
		Action: stdFromTo(decoderTOML, encoderJSON),
	},
	{
		Name:   "toml-to-yaml",
		Usage:  "converts toml to yaml",
		Action: stdFromTo(decoderTOML, encoderYAML),
	},

	// YAML -
	{
		Name:  "yaml-to-csv",
		Usage: `converts json to csv. It must be a list of objects, eg.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "delimiter",
				Aliases: []string{"d"},
				Value:   ",",
			},
		},
		Action: toCsv(decoderYAML),
	},
	{
		Name:   "yaml-to-xml",
		Usage:  "converts yaml to xml (WARNING: works poorly, xml is broken)",
		Action: stdFromTo(decoderYAML, encoderXML),
	},
	{
		Name:   "yaml-to-json",
		Usage:  "converts yaml to json",
		Action: stdFromTo(decoderYAML, encoderJSON),
	},
	{
		Name:   "yaml-to-toml",
		Usage:  "converts toml to yaml",
		Action: stdFromTo(decoderYAML, encoderTOML),
	},

	// XML todo -- needs some special care
	//{
	//	Name:   "xml-to-yaml",
	//	Usage:  "converts xml to yaml (WARNING: works poorly, xml is broken)",
	//	Action: stdFromTo(decoderXML, encoderYAML),
	//},
	//{
	//	Name:   "xml-to-json",
	//	Usage:  "converts xml to json",
	//	Action: stdFromTo(decoderXML, encoderJSON),
	//},
	//{
	//	Name:   "xml-to-toml",
	//	Usage:  "converts xml to toml",
	//	Action: stdFromTo(decoderXML, encoderTOML),
	//},

	// Other

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
