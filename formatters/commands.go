package formatters

import (
	"github.com/urfave/cli/v3"
)

var Commands = []*cli.Command{
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
		Action: formatJSON,
	},
	{
		Name: "xml",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "indent",
				Value: 2,
			},
		},
		Action: formatXML,
	},
	{
		Name:   "lower",
		Action: toLowerCase,
	},
	{
		Name:   "upper",
		Action: toUpperCase,
	},
}
