# Chemtrail Scaling Policies
Scaling policies allow operators to control how Chemtrail scales the Nomad client node worker pool as well as configure critical parameters. Each policy is associated to a Nomad client node class, a defined with the [configuration documentation](https://www.nomadproject.io/docs/configuration/client.html#node_class). Nodes that are not configured with the `node_class` parameter, are considered to be part of the `chemtrail-default` class. 

### Scaling Policy Params
The top level policy parameters control the autoscalers overall actions when interacting with a Nomad client class.

* `Enabled` (bool) - Whether autoscaling assessments will be made for the Nomad client node class. 
* `Class` (string) - The Nomad client class that this scaling policy is for.
* `MinCount` (int) - The minimum number of nodes that should be running in the class pool.
* `MaxCount` (int)  - The maximum number of nodes that should be running in the class pool.
* `ScaleInCount` (int) - The number by which to decrement the node class count by when performing a scaling in action. Currently this can only be `1`.
* `ScaleOutCount` (int) - The number by which to increment the ode class count by when performing a scaling out action.
* `Provider` (string) - The node provider used to perform scaling actions. Currently `aws-autoscaling` is supported.
* `ProviderConfig` (map[string]string) - A key/value map containing configuration to be used when calling the `Provider`.
* `Checks` (map[string]Check) - A map containing the desired checks to perform during an autoscaling evaluation. The key is a free-form user supplied string value, identifying the check. The params of a check are detailed below.

### Scaling Policy Check Params
Multiple checks can be provided per scaling policy. During evaluation runs where two checks decide the opposite action should be triggered, the scale out will always take priority over scale in.

* `Enabled` (bool) - Whether this individual check should be run or not.
* `Resource` (string) - is the Nomad resource to evaluate. This currently supports `cpu` and `memory` as define by the Nomad job [resource stanza](https://www.nomadproject.io/docs/job-specification/resources.html).
* `ComparisonOperator` (string) - The operator used when evaluating a metric value against a threshold. This currently supports `greater-than` and `less-than`.
* `ComparisonPercentage` (float64) - The threshold value compared against the resource allocation percentage to check whether the check action should be triggered.
* `Action` (string) - The action to take if the metric breaks the threshold. This currently supports `scale-in` and `scale-out`.

## Full Example
Below is a full scaling policy example, which contains scale out and scale in checks for CPU and memory metrics.

```json
{
  "Enabled": true,
  "MinCount": 2,
  "MaxCount": 4,
  "ScaleOutCount": 1,
  "ScaleInCount": 1,
  "Provider": "aws-autoscaling",
  "ProviderConfig": {
    "asg-name": "chemtrail-test"
  },
  "Checks": {
    "cpu-in": {
      "Enabled": true,
      "Resource": "cpu",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "cpu-out": {
      "Enabled": true,
      "Resource": "cpu",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    },
    "memory-in": {
      "Enabled": true,
      "Resource": "memory",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "memory-out": {
      "Enabled": true,
      "Resource": "memory",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    }
  }
}
```