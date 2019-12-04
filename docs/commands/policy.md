# Policy CLI

The policy command groups subcommands for interacting with policies. Users can write, read, and list policies in Chemtrail.

## Examples

List all policies:
```bash
$ chemtrail policy list
```

Read a policy for the client node class high-memory:
```bash
$ chemtrail policy read high-memory
```

Create a policy for the client node class high-memory:
```bash
$ chemtrail policy write high-memory policy.json
```

Delete the policy for the client node class high-memory:
```bash
$ chemtrail policy delete high-memory
```

## Usage
```bash
Usage:
  chemtrail policy [flags]
  chemtrail policy [command]

Available Commands:
  delete      Deletes a scaling policy
  init        Creates an example scaling policy
  list        Lists all scaling policies
  read        Details the scaling policy
  write       Uploads a policy from file
```
