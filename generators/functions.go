package generators

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/xid"
	"github.com/urfave/cli/v3"
	"math/big"
	rand2 "math/rand/v2"
	"strconv"
	"strings"
)

func generatePassphrase(ctx context.Context, c *cli.Command) error {
	out := c.Writer

	short := c.Bool("short")
	mix := c.Bool("mix")
	l := 4
	ls := c.Args().Get(0)
	ii, err := strconv.ParseInt(ls, 10, 32)
	if err == nil {
		l = int(ii)
	}

	list := strings.Split(enWordlist, "\n")
	if mix {
		list = append(list, strings.Split(svWordlist, "\n")...)
	}
	rand2.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

	var res []string

	for i := 0; i < l; {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(list))))
		if err != nil {
			return err
		}

		word := list[index.Int64()]

		if short && len(word) > 5 {
			continue
		}

		res = append(res, word)
		i++
	}

	_, err = out.Write([]byte(strings.Join(res, "-")))
	return err
}

func generateUUID(ctx context.Context, c *cli.Command) error {
	out := c.Writer

	version := c.Int("version")
	space, _ := uuid.Parse(c.String("namespace"))
	data := c.String("data")

	var fn func() (uuid.UUID, error)
	switch version {

	case 7:
		fn = uuid.NewV7
	case 6:
		fn = uuid.NewV6
	case 5:
		fn = func() (uuid.UUID, error) {
			return uuid.NewSHA1(space, []byte(data)), nil
		}
	case 4:
		fn = uuid.NewRandom
	case 3:
		fn = func() (uuid.UUID, error) {
			return uuid.NewMD5(space, []byte(data)), nil
		}
	default:
		fn = uuid.NewRandom
	}
	u, err := fn()
	if err != nil {
		return err
	}
	_, err = out.Write([]byte(u.String()))
	return err
}

func generateXID(ctx context.Context, c *cli.Command) error {
	out := c.Writer

	_, err := out.Write([]byte(xid.New().String()))
	return err
}

func generateRandomString(ctx context.Context, c *cli.Command) error {
	out := c.Writer

	l := 16
	ls := strings.Trim(c.Args().Get(0), " \t\"'\\")
	ii, err := strconv.ParseInt(ls, 10, 32)
	if err == nil && ii > 0 {
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
}

func generateRandomBytes(ctx context.Context, c *cli.Command) error {
	out := c.Writer
	l := 16

	ls := strings.Trim(c.Args().Get(0), " \t\"'\\")
	ii, err := strconv.ParseInt(ls, 10, 32)
	if err == nil && ii > 0 {
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
}
