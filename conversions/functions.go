package conversions

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/modfin/henry/mapz"
	"github.com/modfin/henry/slicez"
	"github.com/urfave/cli/v3"
	"io"
	"math/big"
)

func intToString(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	bs, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	if c.Bool("little-endian") {
		bs = slicez.Reverse(bs)
	}

	i := big.NewInt(0)
	i.SetBytes(bs)
	_, err = out.Write([]byte(i.String()))

	return err
}

func stringToInt(ctx context.Context, c *cli.Command) error {
	in := c.Reader
	out := c.Writer

	bs, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	f, _, err := big.ParseFloat(string(bytes.TrimSpace(bs)), 10, 0, big.ToZero)
	if err != nil {
		return err
	}
	i, _ := f.Int(nil)

	b := i.Bytes()
	if len(b) == 0 {
		b = []byte{0}
	}
	if c.Bool("little-endian") {
		b = slicez.Reverse(b)
	}

	_, err = out.Write(b)

	return err
}

type encoder interface {
	Encode(v any) error
}
type decoder interface {
	Decode(v any) error
}

func stdFromTo(decode func(r io.Reader) decoder, encode func(w io.Writer) encoder) cli.ActionFunc {
	return func(ctx context.Context, command *cli.Command) error {
		in := command.Reader
		out := command.Writer

		var item any

		err := decode(in).Decode(&item)
		if err != nil {
			return err
		}
		return encode(out).Encode(item)
	}
}

func csvTo(toEnc func(w io.Writer) encoder) func(ctx context.Context, c *cli.Command) error {
	return func(ctx context.Context, c *cli.Command) error {
		var err error
		enc := toEnc(c.Writer)
		in := c.Reader

		reader := csv.NewReader(in)

		switch c.String("delimiter") {
		case "\\t":
			reader.Comma = '\t'
		case "\\n":
			reader.Comma = '\n'
		case "":
			reader.Comma = ','
		default:
			reader.Comma = rune(c.String("delimiter")[0])
		}
		reader.LazyQuotes = true

		var headers []string
		if c.Bool("with-headers") {
			headers, err = reader.Read()
			if err != nil {
				return err
			}
		}

		getName := func(col int) string {
			if len(headers) > col {
				return headers[col]
			}
			return fmt.Sprintf("col_%d", col)
		}

		rows := []map[string]string{}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}

			row := map[string]string{}
			for i, v := range record {
				row[getName(i)] = v
			}

			rows = append(rows, row)
		}

		err = enc.Encode(rows)
		if err != nil {
			return err
		}

		return nil
	}

}

func toCsv(decode func(r io.Reader) decoder) func(ctx context.Context, c *cli.Command) error {
	return func(ctx context.Context, c *cli.Command) error {
		in := c.Reader
		out := c.Writer

		var items []map[string]interface{}
		err := decode(in).Decode(&items)
		if err != nil {
			return fmt.Errorf("failed to decode: %s", err)
		}

		if len(items) == 0 {
			return nil
		}

		writer := csv.NewWriter(out)

		switch c.String("delimiter") {
		case "\\t":
			writer.Comma = '\t'
		case "\\n":
			writer.Comma = '\n'
		case "":
			writer.Comma = ','
		default:
			writer.Comma = rune(c.String("delimiter")[0])
		}

		var headers []string
		headers = slicez.Uniq(slicez.FlatMap(items, func(item map[string]interface{}) []string {
			return mapz.Keys(item)
		}))

		err = writer.Write(headers)
		if err != nil {
			return fmt.Errorf("failed to write headers: %s", err)
		}

		for _, item := range items {
			var rec []string
			for _, v := range headers {
				if item[v] == nil {
					rec = append(rec, "")
					continue
				}
				rec = append(rec, fmt.Sprintf("%v", item[v]))
			}
			err = writer.Write(rec)
			if err != nil {
				return fmt.Errorf("failed to write record: %s", err)
			}
		}

		writer.Flush()

		return nil
	}

}
