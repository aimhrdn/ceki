# CEKI

CEKI (check key) is cli tool for checking validity of your AWS KEY. 

# Compile tools

```
## Clone source code
git clone https://github.com/mahirrudin/ceki
cd ceki

## windows
env GOOS=windows GOARCH=amd64 go build -o ceki.exe

## winux
env GOOS=linux GOARCH=amd64 go build -o ceki

```

# Usage of tools

```
A tool to validate and check AWS keys

Usage:
  ceki [command]

Available Commands:
  check       Check AWS S3 access
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  validate    Validate AWS access keys

Flags:
  -h, --help   help for ceki

Use "ceki [command] --help" for more information about a command.
```