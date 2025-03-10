package generators

import (
	_ "embed"
	"github.com/urfave/cli/v3"
	"time"
)

//go:embed wordlist_en
var enWordlist string

//go:embed wordlist_sv
var svWordlist string

var Commands = []*cli.Command{
	{
		Name: "pass",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "short",
				Usage: "returns shorter words",
			},
			&cli.BoolFlag{
				Name:  "mix",
				Usage: "mixes swedish and english",
			},
		},
		Usage:     "generates a random passphrase",
		Action:    generatePassphrase,
		ArgsUsage: "[words] (default 4)",
	},
	{
		Name:  "uuid",
		Usage: "generate a uuid",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "version",
				Value:   4,
				Usage:   "specify the uuid version. Supported 3-7 inclusive",
				Aliases: []string{"v"},
			},

			&cli.StringFlag{
				Name:    "namespace",
				Usage:   "uuid used as namespace. Used with version 3 and 5",
				Aliases: []string{"ns"},
			},
			&cli.StringFlag{
				Name:    "data",
				Usage:   "data used to generate uuid. Used with version 3 and 5",
				Aliases: []string{"d"},
			},
		},
		Action: generateUUID,
	},
	{
		Name:  "xid",
		Usage: "generate a xid, https://github.com/rs/xid",

		Flags: []cli.Flag{
			&cli.TimestampFlag{
				Name: "time",
				Config: cli.TimestampConfig{
					Timezone: time.Local,
					Layouts: []string{
						time.DateTime,
						time.RFC3339,
					},
				},
			},
			&cli.IntFlag{
				Name: "counter",
			},
			&cli.IntFlag{
				Name: "machine",
			},
			&cli.IntFlag{
				Name: "pid",
			},
		},
		Action: generateXID,
	},
	{
		Name:      "string",
		Usage:     "generate a random string",
		ArgsUsage: "[length] (default 16)",
		Action:    generateRandomString,
	},
	{
		Name:      "bytes",
		Usage:     "generate random bytes",
		ArgsUsage: "[length] (default 16)",
		Action:    generateRandomBytes,
	},
}
