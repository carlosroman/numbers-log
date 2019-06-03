Go number logger
==================

[![Go Report Card](https://goreportcard.com/badge/github.com/carlosroman/numbers-log)](https://goreportcard.com/report/github.com/carlosroman/numbers-log)
[![CircleCI](https://circleci.com/gh/carlosroman/numbers-log.svg?style=svg)](https://circleci.com/gh/carlosroman/numbers-log)
[![codecov](https://codecov.io/gh/carlosroman/numbers-log/branch/master/graph/badge.svg)](https://codecov.io/gh/carlosroman/numbers-log)



## Install

The project requires the following:
* Golang (1.11+)
* GNU Make (optional)

Clone the project:

```
$ git clone https://github.com/carlosroman/numbers-log.git
```

## Build

To build run:
```
$ make build
```
This will produce an executable called `server` in a folder called `target`.

## Run

The quickest way to run it is to either build the exe and then `./target/server` or :

```
$ go run cmd/server/main.go
```

## Tests

To run the tests just run:

```
$ make test
```

## License

MIT.
