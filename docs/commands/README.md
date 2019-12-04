# Chemtrail Commands (CLI)

In addition to a verbose HTTP API, Chemtrail features a command-line interface that wraps common functionality and formats output.

To get help, run:
```bash
$ chemtrail -h
```

To get help for a subcommand, run:
```bash
$ chemtrail <subcommand> -h
```

## CLI Command Structure

There are a number of command and subcommand options available. Construct your Chemtrail CLI command such that the command options precede its and arguments if any:

```bash
chemtrail <command> <subcommand> [options] [args]
```

## General Options

The following options are available to all Chemtrail CLI commands and help set client connection variables.

* `--addr` (string: "http://127.0.0.1:8000") - The HTTP(S) address of the Chemtrail server
* `--client-ca-path` (string: "") - Path to a PEM encoded CA cert file to use to verify the Chemtrail server SSL certificate.
* `--client-cert-key-path` (string: "") - Path to an unencrypted PEM encoded private key matching the client certificate
* `--client-cert-path string` (string: "") - Path to a PEM encoded client certificate for TLS authentication to the Chemtrail server

## Exit Codes

The Chemtrail CLI aims to be consistent and well-behaved using [sysexits](https://github.com/sean-/sysexits) for CLI exit codes.

* exit code `64`: represents local errors such as incorrect flags, failed validation or an incorrect number of passed arguments.
* exit code `70`: represents an internal failure such as API failures.
