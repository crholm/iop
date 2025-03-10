package main

import (
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/crholm/iop/conversions"
	"github.com/crholm/iop/decoders"
	"github.com/crholm/iop/encoders"
	"github.com/crholm/iop/formatters"
	"github.com/crholm/iop/generators"
	"github.com/urfave/cli/v3"
	"io"
	"os"
	"os/exec"
	"sync"
)

func main() {

	var commands [][]string
	var command []string
	for _, a := range os.Args[1:] {
		if a == "--" {
			commands = append(commands, command)
			command = nil
			continue
		}
		command = append(command, a)
	}
	commands = append(commands, command)

	if len(commands) > 1 {
		multiSpawn(commands)
	}

	err := createApp().Run(context.Background(), append([]string{"/proc/self/exe"}, commands[0]...))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "got err", err)
		os.Exit(1)
	}

}

func multiSpawn(commands [][]string) {
	var wg sync.WaitGroup
	wg.Add(len(commands))

	cmds := make([]*exec.Cmd, len(commands))
	for i, args := range commands {
		cmds[i] = exec.Command("/proc/self/exe", args...)
	}

	for i, cmd := range cmds {
		if i == 0 {
			cmd.Stdin = os.Stdin
		}
		if i > 0 {
			cmd.Stdin, _ = cmds[i-1].StdoutPipe()
		}
		if i == len(commands)-1 {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr

		go func(cmd *exec.Cmd) {
			defer wg.Done()
			defer cmd.Wait() // Very important to let all io loops finish writing to stdout and stderr
			err := cmd.Start()
			if err != nil {
				os.Stderr.Write([]byte(err.Error()))
			}

		}(cmd)
	}

	wg.Wait()

	os.Exit(0)
}

func createApp() *cli.Command {

	app := &cli.Command{
		Name:      "iop",
		Usage:     "a tool for converting and formatting things from std in to std out",
		UsageText: "You can use -- as piping between commands, eg. echo 124 | iop conv string-to-int -- encode hex -- clip copy",
		Commands: []*cli.Command{

			{
				Name:  "copy",
				Usage: "puts things from std in onto the clipboard",
				Action: func(ctx context.Context, c *cli.Command) error {
					in := c.Reader
					b, err := io.ReadAll(in)
					if err != nil {
						return err
					}
					return clipboard.WriteAll(string(b))
				},
			},
			{
				Name:  "paste",
				Usage: "puts things in clipboard onto std out",
				Action: func(ctx context.Context, c *cli.Command) error {
					out := c.Writer

					str, err := clipboard.ReadAll()
					if err != nil {
						return err
					}

					_, err = out.Write([]byte(str))
					return err
				},
			},

			{
				Name:     "decode",
				Aliases:  []string{"dec"},
				Usage:    "decode std from something",
				Commands: decoders.Commands,
			},
			{
				Name:     "encode",
				Aliases:  []string{"enc"},
				Usage:    "encode std to something",
				Commands: encoders.Commands,
			},
			{
				Name:     "fmt",
				Aliases:  []string{"format"},
				Usage:    "format something from std in",
				Commands: formatters.Commands,
			},
			{
				Name:     "gen",
				Aliases:  []string{"generate"},
				Usage:    "generate something",
				Commands: generators.Commands,
			},

			{
				Name:     "conv",
				Aliases:  []string{"convert"},
				Usage:    "convert something",
				Commands: conversions.Commands,
			},
		},
	}
	return app
}
