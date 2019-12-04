# Scale API

## Scale Out Client Node Class Group

This endpoint can be used to scale a Nomad client node class out, therefore increasing its count.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/v1/scale/out/:client_class`              | `201 application/binary` |

#### Parameters

* `:client_class` (string: required) - Specifies the client node class to scale out

### Sample Request

```
$ curl \
    --request POST \
    http://127.0.0.1:8000/v1/scale/out/high-memory
```

### Sample Response

```json
{
  "ID": "036e4bd6-8f7d-4a8c-bf90-790790bbdc2a"
}
```

## Scale In Client Node Class Group

This endpoint can be used to scale a Nomad client node class in, therefore decreasing its count.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `POST`    | `/v1/scale/in/:client_class`              | `201 application/binary` |

#### Parameters

* `:client_class` (string: required) - Specifies the client node class to scale in

### Sample Request

```
$ curl \
    --request POST \
    http://127.0.0.1:8000/v1/scale/in/high-memory
```

### Sample Response

```json
{
  "ID": "036e4bd6-8f7d-4a8c-bf90-790790bbdc2a"
}
```

## List Scaling Events

This endpoint can be used to list the recent scaling events.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/scale/status`              | `200 application/binary` |

### Sample Request

```
$ curl \
    --request GET \
    http://127.0.0.1:8000/v1/scale/status
```

### Sample Response

```json
{
  "6114c740-5656-4313-bd59-4ac1baa6c7b2": {
    "Events": [
      {
        "Timestamp": 1674947513737560000,
        "Message": "scaling activity has started",
        "Source": "chemtrail"
      },
      {
        "Timestamp": 1674947513748519000,
        "Message": "all allocations on node \"d832b8c2-1b8d-72ce-8f2b-8b2610d0aaf9\" have stopped",
        "Source": "nomad"
      },
      {
        "Timestamp": 1674947513819148000,
        "Message": "drain complete for node d832b8c2-1b8d-72ce-8f2b-8b2610d0aaf9",
        "Source": "nomad"
      },
      {
        "Timestamp": 1674947513821821000,
        "Message": "scaling activity has successfully completed",
        "Source": "chemtrail"
      }
    ],
    "Direction": "in",
    "LastUpdate": 1674947513821821000,
    "Status": "success",
    "Provider": "aws-autoscaling",
    "ProviderCfg": {
      "asg-name": "chemtrail-test"
    }
  },
  "8aa2d024-0dea-45c4-875e-78c89acf9d93": {
    "Events": [
      {
        "Timestamp": 1674947693736639000,
        "Message": "scaling activity has started",
        "Source": "chemtrail"
      },
      {
        "Timestamp": 1674947693745198000,
        "Message": "all allocations on node \"d832b8c2-1b8d-72ce-8f2b-8b2610d0aaf9\" have stopped",
        "Source": "nomad"
      },
      {
        "Timestamp": 1674947693805382000,
        "Message": "drain complete for node d832b8c2-1b8d-72ce-8f2b-8b2610d0aaf9",
        "Source": "nomad"
      },
      {
        "Timestamp": 1674947693808511000,
        "Message": "scaling activity has successfully completed",
        "Source": "chemtrail"
      }
    ],
    "Direction": "in",
    "LastUpdate": 1674947693808511000,
    "Status": "success",
    "Provider": "aws-autoscaling",
    "ProviderCfg": {
      "asg-name": "chemtrail-test"
    }
  }
}
```

## Read Scaling Event

This endpoint can be used to query a scaling event.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/scale/status/:id`              | `200 application/binary` |

### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/scale/status/8aa2d024-0dea-45c4-875e-78c89acf9d93
```

### Sample Response

```json
{
  "Events": [
    {
      "Timestamp": 1674947693736639000,
      "Message": "scaling activity has started",
      "Source": "chemtrail"
    },
    {
      "Timestamp": 1674947693745198000,
      "Message": "all allocations on node \"d832b8c2-1b8d-72ce-8f2b-8b2610d0aaf9\" have stopped",
      "Source": "nomad"
    },
    {
      "Timestamp": 1674947693805382000,
      "Message": "drain complete for node d832b8c2-1b8d-72ce-8f2b-8b2610d0aaf9",
      "Source": "nomad"
    },
    {
      "Timestamp": 1574947693808511000,
      "Message": "scaling activity has successfully completed",
      "Source": "chemtrail"
    }
  ],
  "Direction": "in",
  "LastUpdate": 1674947693808511000,
  "Status": "success",
  "Provider": "aws-autoscaling",
  "ProviderCfg": {
    "asg-name": "chemtrail-test"
  }
}
```
