# Health-check for `cardano-node`

A simple command-line utility to locally check whether a running instance of `cardano-node` is healthy. Local
refers here to the fact that the tool isn't using any remote information in order to come to a conclusion.

## Build

```bash
$ go mod vendor
$ go build -o healthcheck
```

## Run

```
Usage:
        ./healthcheck [flags]
Flags:
  -config-file string
        path to the configuration file of cardano-node
  -genesis-file string
        path to the genesis file of cardano-node
  -max-time-since-last-block duration
        threshold for duration between now and the creation date of the most recently received block (default 10m0s)
```

## Contributions

Feel free to submit tickets and pull requests. This repository is maintained by:

* [Kevin Haller](mailto:kevin.haller@blockblu.io) (Operator of the [SOBIT](https://staking.outofbits.com/) stake pool)
