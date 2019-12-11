# Chemtrail Documentation

Chemtrail is a client node scaler for [HashiCorp Nomad](https://www.nomadproject.io/) allowing for dynamic scaling of class worker pools based on demand.

### Key Features
* __Scale node worker pool based on allocated demand:__ The autoscaler uses Nomad resource allocation metrics to dynamically scale client class worker pools. This ensures both capacity availability to meet demand, and cost efficiency.
* __Operator friendly:__ Chemtrail is designed to be easy to operate but flexible. Scaling state provides detailed insights into the actions undertaken during an autoscaling event.
* __Easily extensible to scale cloud or physical host providers:__ The provider interface is simple and concise, allowing for easy extension to support your desired cloud of physical server provider.

## Table of contents
1. [API](./api) documentation.
1. [CLI](./commands) documentation.
1. [Chemtrail server](./configuration) configuration documentation.
1. [Guides](./guides) documentation.
