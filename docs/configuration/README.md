# Chemtrail Configuration

The Chemtrail server can be configured by supplying either CLI flags or using environment variables.

## Parameters

* `--autoscaler-enabled` (bool: false) - Enable the internal autoscaling engine
* `--autoscaler-evaluation-interval` (int: 180) - The time period in seconds between autoscaling evaluation runs.
* `--autoscaler-num-threads` (int: 3) - Specifies the number of parallel autoscaler threads to run.
* `--bind-addr` (string: "127.0.0.1") - The HTTP server address to bind to.
* `--bind-port` (uint16: 8000) - The HTTP server port to bind to.
* `--log-enable-dev` (bool: false) - Log with file:line of the caller.
* `--log-format` (string: "auto") - Specify the log format ("auto", "zerolog" or "human").
* `--log-level` (string: "info") - Change the level used for logging.
* `--log-use-color` (bool: false) - Use ANSI colors in logging output.
* `--provider-aws-asg-enabled` (bool: false) - Enable the AWS AutoScaling Group client provider.
* `--provider-noop-enabled` (bool: true) - Enable the NoOp client provider.
* `--storage-consul-enabled` (bool: false) - Use Consul as the storage backend for state.
* `--storage-consul-path` (string: "chemtrail/") - The Consul KV path that will be used to store policies and state.
* `--telemetry-statsd-address` (string: "") - Specifies the address of a statsd server to forward metrics to.
* `--telemetry-statsite-address` (string: "") - Specifies the address of a statsite server to forward metrics data to.
* `--tls-cert-key-path` (string: "") - Path to the TLS certificate key for the Chemtrail server.
* `--tls-cert-path` (string: "") - Path to the TLS certificate for the Chemtrail server.

### Environment Variables

When specifying environment variables, the CLI flag should be converted like follows:
* `--bind-addr` becomes `CHEMTRAIL_BIND_ADDR`

## Client Parameters

Nomad and Consul clients can be configured using the native environment variables which are available through the HashiCorp SDKs. Using these keeps the setup simple and consistent.

### Nomad Client Parameters

The Nomad client environment variables documentation can be found on the [Nomad general options](https://github.com/hashicorp/nomad/blob/22fd62753510a4a41c1b8f1d117ea1a90b48df06/website/source/docs/commands/_general_options.html.md) GitHub document. For ease of use this document is reproduced below:

* `NOMAD_ADDR` (string: "http://127.0.0.1:4646") - The address of the Nomad server.
* `NOMAD_REGION` (string: "") - The region of the Nomad servers to forward commands to.
* `NOMAD_NAMESPACE` (string "default") - The target namespace for queries and actions bound to a namespace.
* `NOMAD_CACERT` (string: "") - Path to a PEM encoded CA cert file to use to verify the Nomad server SSL certificate.
* `NOMAD_CAPATH` (string: "") - Path to a directory of PEM encoded CA cert files to verify the Nomad server SSL certificate.
* `NOMAD_CLIENT_CERT` (string: "") - Path to a PEM encoded client certificate for TLS authentication to the Nomad server.
* `NOMAD_CLIENT_KEY` (string: "") - Path to an unencrypted PEM encoded private key matching the client certificate.
* `NOMAD_SKIP_VERIFY` (bool: false) - Do not verify TLS certificate.
* `NOMAD_TOKEN` (string: "") - The SecretID of an ACL token to use to authenticate API requests with.

### Consul Client Parameters

The Consul client environment variables documentation can be found on the [Consul commands page](https://www.consul.io/docs/commands/index.html#environment-variables). For ease of use this document is reproduced below:

* `CONSUL_HTTP_ADDR` (string: "127.0.0.1:8500") - This is the HTTP API address to the local Consul agent (not the remote server) specified as a URI with optional scheme.
* `CONSUL_HTTP_TOKEN` (string: "") - This is the API access token required when access control lists (ACLs) are enabled.
* `CONSUL_HTTP_AUTH` (string: "") - This specifies HTTP Basic access credentials as a username:password pair.
* `CONSUL_HTTP_SSL` (bool: false) - This is a boolean value that enables the HTTPS URI scheme and SSL connections to the HTTP API.
* `CONSUL_HTTP_SSL_VERIFY` (bool: true) - This is a boolean value to specify SSL certificate verification.
* `CONSUL_CACERT` (string: "") - Path to a CA file to use for TLS when communicating with Consul.
* `CONSUL_CAPATH` (string: "") - Path to a directory of CA certificates to use for TLS when communicating with Consul.
* `CONSUL_CLIENT_CERT` (string: "") - Path to a client cert file to use for TLS.
* `CONSUL_CLIENT_KEY` (string: "") - Path to a client key file to use for TLS.
* `CONSUL_TLS_SERVER_NAME` (string: "") - The server name to use as the SNI host when connecting via TLS.

## Provider Configuration

In order to make scaling requests to backend providers, configuration is required providing authentication amongst others. Below are specific details of the minimum requirement for each provider, and links if available to more in-depth documentation.

### NoOp

The NoOp provider is ideal when first introducing Chemtrail to your environment, or when making configuration changes. The NoOp provider will not enact scaling changes, but will log at INFO level intended behaviour and predicted actions.

### Amazon Web Services Auto Scaling Groups

Chemtrail is using the official AWS GO SDK so you can authenticate to AWS with static credentials, instance role or environment variables, details of which can be found on the [AWS site](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html). 

For ease the most commonly used and required environment variables can be found below:

* `AWS_ACCESS_KEY_ID` - Specifies an AWS access key associated with an IAM user or role.
* `AWS_SECRET_ACCESS_KEY` - Specifies the secret key associated with the access key. This is essentially the "password" for the access key.
* `AWS_DEFAULT_REGION` - Specifies the AWS Region to send the request to.

The required IAM permissions are :

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Chemtrail",
            "Effect": "Allow",
            "Action": [
                "ec2:TerminateInstances",
                "autoscaling:UpdateAutoScalingGroup",
                "autoscaling:DetachInstances",
                "autoscaling:DescribeAutoScalingGroups"
            ],
            "Resource": "*"
        }
    ]
}
```
