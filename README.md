

> A small cli tool for some daily needs. 

# Installation
```go install https://github.com/crholm/iop```

# Usage

## iop

```text
NAME:
   iop - a tool for converting and formatting things from std in to std out

USAGE:
   iop [global options] command [command options] [arguments...]

COMMANDS:
   decode   decode std from something
   encode   encode std to something
   fmt      format something from std in
   rand     generate something random
   help, h  Shows a list of commands or help for one command

```

## Decode
```text
NAME:
   iop decode - decode std from something

USAGE:
   iop decode command [command options] [arguments...]

COMMANDS:
   url          decodes a url query string
   b64, base64  decodes a base64 string
   b32, base32  decodes a base32 string
   hex          decodes a hex string
   jwt          decodes a jwt token
   help, h      Shows a list of commands or help for one command
```

## Encode
```text
NAME:
   iop encode - encode std to something

USAGE:
   iop encode command [command options] [arguments...]

COMMANDS:
   url          url query encodes a data
   b64, base64  base64 encodes a data
   b32, base32  base32 encodes a data
   hex          hex encodes a data
   help, h      Shows a list of commands or help for one command
```

## fmt
```text
NAME:
   iop fmt - format something from std in

USAGE:
   iop fmt command [command options] [arguments...]

COMMANDS:
   json     
   xml      
   help, h  Shows a list of commands or help for one command
```

## rand
```text
NAME:
   iop rand - generate something random

USAGE:
   iop rand command [command options] [arguments...]

COMMANDS:
   uuid     generate a random v4 uuid
   string   generate a random string
   bytes    generate random bytes
   help, h  Shows a list of commands or help for one command

```