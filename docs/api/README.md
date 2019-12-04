# HTTP API

The Chemtrail HTTP API gives you full access to a Chemtrail server via HTTP. Every aspect of Chemtrail can be controlled via this API.

All API routes are prefixed with /v1/, which is the current API version.

## Table of contents
1. [Policy API](./policy.md) documentation.
1. [Scale API](./scale.md) documentation.
1. [System API](./system.md) documentation.

## HTTP Status Codes
The following HTTP status codes are used throughout the API. Chemtrail tries to adhere to these whenever possible.

* `200` - Success with data.
* `201` - Success created without return content.
* `404` - Not found.
* `422` - Unprocessable request. An error where the supplied payload or query params are incorrect.
* `500` - Internal server error. An internal error has occurred, try again later.
