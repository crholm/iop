

# IOP - Input/Output Processor

IOP is a versatile command-line utility for converting, encoding, decoding, formatting, and generating data. It's designed to process data from standard input to standard output, making it easy to integrate into command pipelines.

## Installation

```bash
go install github.com/crholm/iop
```

## Basic Usage

IOP follows a simple pattern: it reads from stdin, processes the data according to the specified command, and writes to stdout.

```bash
echo "Hello World" | iop encode base64
```

## Key Features

### Clipboard Integration

Copy to and paste from the clipboard:

```bash
# Copy text to clipboard
cat file.txt | iop copy

# Paste from clipboard to stdout
iop paste > file.txt
```

### Encoding & Decoding

Support for various encoding formats:

```bash
# Encode to base64
echo "Hello World" | iop encode base64

# Decode from hex
echo "48656c6c6f20576f726c64" | iop decode hex
```

### Formatting

Format and prettify various data formats:

```bash
# Format JSON with 2-space indentation
cat data.json | iop fmt json --indent 2

# Format XML
cat data.xml | iop fmt xml
```

### Type Conversion

Convert between data types:

```bash
# Convert string to integer
echo "123" | iop conv string-to-int

# Convert integer to string
echo -ne "\x7B" | iop conv int-to-string
```

### Data Generation

Generate various types of data:

```bash
# Generate a UUID
iop gen uuid

# Generate a random string
iop gen string 32

# Generate a passphrase
iop gen passphrase
```

## Pipeline Chaining

One of the most powerful features of IOP is the ability to chain commands using `--`:

```bash
# Convert string to int, encode as hex, and copy to clipboard
echo "123" | iop conv string-to-int -- encode hex -- copy
```

## Command Reference

### Encoding Commands
- `encode base64` - Encode data to base64
- `encode base32` - Encode data to base32
- `encode hex` - Encode data to hexadecimal
- `encode binary` - Encode data to binary
- `encode url` - Encode data for URLs (query params)

### Decoding Commands
- `decode base64` - Decode base64 data
- `decode base32` - Decode base32 data
- `decode hex` - Decode hexadecimal data
- `decode binary` - Decode binary data
- `decode url` - Decode URL-encoded data (query params)

### Format Commands
- `fmt json` - Format JSON data
- `fmt xml` - Format XML data
- `fmt lower` - Convert text to lowercase
- `fmt upper` - Convert text to uppercase

### Conversion Commands
- `conv string-to-int` - Convert string to integer
- `conv int-to-string` - Convert integer to string
- `conv csv-to-json` - Convert CSV to JSON
- `conv csv-to-yaml` - Convert CSV to YAML
- `conv csv-to-xml` - Convert CSV to xml

### Generator Commands
- `gen uuid [--version N]` - Generate a UUID (versions 3-7)
- `gen xid` - Generate an XID
- `gen random-string [length]` - Generate a random string
- `gen random-bytes [length]` - Generate random bytes
- `gen passphrase [count] [--short] [--mix]` - Generate a passphrase

### Clipboard Commands
- `copy` - Copy stdin to clipboard
- `paste` - Paste clipboard to stdout

## Examples

```bash
# Generate a UUID and encode it as base64
iop gen uuid -- encode base64

# Format JSON from clipboard and copy it back
iop paste -- fmt json --indent 2 -- copy

# Generate a random string and convert to uppercase
iop gen random-string 16 -- fmt upper

# Decode a base64 string and format as JSON
echo "eyJuYW1lIjoiSm9obiJ9" | iop decode base64 -- fmt json
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.