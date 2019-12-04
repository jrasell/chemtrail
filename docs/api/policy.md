# Policy API

The policy API allows interaction with scaling policies registered with Chemtrail.

## List Scaling Policies

This endpoint lists all known scaling policies in the system registered with Chemtrail.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/policies`              | `200 application/binary` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/policies
```

### Sample Response

```json
{
  "default": {
    "Enabled": true,
    "Class": "default",
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
        "ScaleResource": "cpu",
        "ComparisonOperator": "less-than",
        "ComparisonPercentage": 25,
        "Action": "scale-in"
      },
      "cpu-out": {
        "Enabled": true,
        "ScaleResource": "cpu",
        "ComparisonOperator": "greater-than",
        "ComparisonPercentage": 80,
        "Action": "scale-out"
      },
      "memory-in": {
        "Enabled": true,
        "ScaleResource": "memory",
        "ComparisonOperator": "less-than",
        "ComparisonPercentage": 25,
        "Action": "scale-in"
      },
      "memory-out": {
        "Enabled": true,
        "ScaleResource": "memory",
        "ComparisonOperator": "greater-than",
        "ComparisonPercentage": 80,
        "Action": "scale-out"
      }
    }
  }
}
```

## Read A Scaling Policy

This endpoint is used to read the scaling policy for a client node class.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/policy/:client_class`              | `200 application/binary` |

#### Parameters

* `:client_class` (string: required) - Specifies the client node class to view the scaling policy of

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/policy/high-memory
```

### Sample Response

```json
{
  "Enabled": true,
  "Class": "high-memory",
  "MinCount": 2,
  "MaxCount": 4,
  "ScaleOutCount": 1,
  "ScaleInCount": 1,
  "Provider": "aws-autoscaling",
  "ProviderConfig": {
    "asg-name": "chemtrail-high-memory"
  },
  "Checks": {
    "cpu-in": {
      "Enabled": true,
      "ScaleResource": "cpu",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "cpu-out": {
      "Enabled": true,
      "ScaleResource": "cpu",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    },
    "memory-in": {
      "Enabled": true,
      "ScaleResource": "memory",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "memory-out": {
      "Enabled": true,
      "ScaleResource": "memory",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    }
  }
}
```

## Create/Update A Scaling Policy

This endpoint can be used to create or update the scaling policy.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/v1/policy/:client_class`              | `200 application/binary` |

#### Parameters

* `:client_class` (string: required) - Specifies the client node class and is specified as part of the path.

### Sample Payload

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
      "ScaleResource": "cpu",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "cpu-out": {
      "Enabled": true,
      "ScaleResource": "cpu",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    },
    "memory-in": {
      "Enabled": true,
      "ScaleResource": "memory",
      "ComparisonOperator": "less-than",
      "ComparisonPercentage": 25,
      "Action": "scale-in"
    },
    "memory-out": {
      "Enabled": true,
      "ScaleResource": "memory",
      "ComparisonOperator": "greater-than",
      "ComparisonPercentage": 80,
      "Action": "scale-out"
    }
  }
}
```

### Sample Request

```
$ curl \
    --request PUT \
    --data @payload.json \
    http://127.0.0.1:8000/v1/policy/general-compute
```

## Delete A Scaling Policy

This endpoint can be used to delete the scaling policy for a client node class.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `DELETE`    | `/v1/policy/:client_class`              | `200 application/binary` |

#### Parameters

* `:client_class` (string: required) - Specifies the client node class and is specified as part of the path.

### Sample Request

```
$ curl \
    --request DELETE \
    http://127.0.0.1:8000/v1/policy/general-compute
```
