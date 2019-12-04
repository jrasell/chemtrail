# System API

## Get Server Health

This endpoint can be used to query the Chemtrail server health status.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/system/health`              | `200 application/binary` |


### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/system/health
```

### Sample Response

```json
{
  "status": "ok"
}
```

## Get Server Metrics

This endpoint can be used to query the Chemtrail server for its latest telemetry data.

| Method   | Path                         |
| :--------------------------- | :--------------------- |
| `GET`    | `/v1/system/metrics`              | `200 application/binary` |


### Sample Request

```
$ curl \
    http://127.0.0.1:8000/v1/system/metrics
```

### Sample Response

```json
{
  "Timestamp": "2019-11-28 19:33:40 +0000 UTC",
  "Gauges": [
    {
      "Name": "chemtrail.runtime.alloc_bytes",
      "Value": 2496696,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.free_count",
      "Value": 63885,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.heap_objects",
      "Value": 22429,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.malloc_count",
      "Value": 86314,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.num_goroutines",
      "Value": 13,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.sys_bytes",
      "Value": 72810744,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.total_gc_pause_ns",
      "Value": 111811,
      "Labels": {}
    },
    {
      "Name": "chemtrail.runtime.total_gc_runs",
      "Value": 2,
      "Labels": {}
    }
  ],
  "Points": [],
  "Counters": [],
  "Samples": []
}
```
