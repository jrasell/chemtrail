# Chemtrail

[![Build Status](https://travis-ci.org/jrasell/chemtrail.svg?branch=master)](https://travis-ci.org/jrasell/chemtrail) [![Go Report Card](https://goreportcard.com/badge/github.com/jrasell/chemtrail)](https://goreportcard.com/report/github.com/jrasell/chemtrail) [![GoDoc](https://godoc.org/github.com/jrasell/chemtrail?status.svg)](https://godoc.org/github.com/jrasell/chemtrail)

Chemtrail is a client scaler for [HashiCorp Nomad](https://www.nomadproject.io/) allowing for dynamic and safe scaling of the client workerpool based on demand.

## Download & Install

* The Chemtrail binary can be downloaded from the [GitHub releases page](https://github.com/jrasell/chemtrail/releases) using `curl -L https://github.com/jrasell/chemtrail/releases/download/v0.0.1/chemtrail_0.0.1_linux_amd64 -o chemtrail`

* A docker image can be found on [Docker Hub](https://hub.docker.com/r/jrasell/chemtrail/), the latest version can be downloaded using `docker pull jrasell/chemtrail`.

* Chemtrail can be built from source by cloning the repository `git clone github.com/jrasell/chemtrail.git` and then using the `make build` command. 

## Documentation

Please refer to the [documentation](./docs) directory for guides to help with deploying and using Chemtrail in your Nomad setup.

## Contributing

Contributions to Chemtrail are very welcome! Please reach out if you have any questions.
