# Scale CLI

The scale command groups subcommands for actioning and detailing scaling requests. Users can scale in or out selected client node classes or view past events.

## Examples

Scale out client node class high-memory:
```bash
$ chemtrail scale out high-memory
```

Scale in client node class high-memory:
```bash
$ chemtrail scale in high-memory
```

List all the scaling events currently held with the Chemtrail storage backend:
```bash
$ chemtrail scale status
```

Read details about the scaling event with id `f7476465-4d6e-c0de-26d0-e383c49be941`:
```
$ chemtrail scale status f7476465-4d6e-c0de-26d0-e383c49be941
```

## Usage
```bash
Usage:
  chemtrail scale [flags]
  chemtrail scale [command]

Available Commands:
  in          Perform scaling in actions on Nomad clients
  out         Perform scaling out actions on Nomad clients
  status      Display the status output for scaling activities
```
