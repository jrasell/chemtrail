# System CLI

The system command groups subcommands gaining insights and information about a running Sherpa server.

## Examples

Detail the health of the server:
```bash
$ chemtrail system health
```

Output the latest metric data points for the running server:
```bash
$ chemtrail system metrics
```

## Usage
```bash
Usage:
  chemtrail system [flags]
  chemtrail system [command]

Available Commands:
  health      Retrieve health information of a Chemtrail server
  metrics     Retrieve metrics from a Chemtrail server
```
