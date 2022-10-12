# Unit-test
Golang unit testing examples.

Please refer to `client`, `server`, `repository`, `eventbus` and `repository/storage` for packages containing various unit tests. 

## Prerequisites

1. Have a running postgres server on localhost, port 5432.
1. Install `psql`.

## Getting started

Run `make db` to instantiate the database used to test the `repository` package.

Run `make test` to run all tests.
