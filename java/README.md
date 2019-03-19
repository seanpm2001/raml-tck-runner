## About

Simple test of few RAML Java parsers. Tests simply try to parse a set of examples and report if parser returned an error.

Running tests produces JSON reports.

A fine collection of RAML files can be composed each containing one/few RAML features to test those in isolation.

Uses [raml-tck](https://github.com/raml-org/raml-tck/tree/master/tests/raml-1.0) as a source of RAML for tests.

NOTE: If file name contains "invalid" parsing of it is expected to fail.

## Install

```sh
$ git clone git@github.com:raml-org/raml-tck-runner.git
$ cd raml-tck-runner/java
$ ./gradlew build
```

## Run

```sh
$ ./gradlew run --args='--parser PARSER_NAME --outdir ./reports/json --branch develop'
```

## Options

Parser:
```sh
$ ./gradlew run --args='--parser webapi-parser/raml-java-parser --branch develop'
```

Output JSON report directory (defaults to `./`):
```sh
$ ./gradlew run --args='--parser webapi-parser --outdir ../reports/json --branch develop'
```

raml-tck branch to load RAML files from:
```sh
$ ./gradlew run --args='--branch develop --parser webapi-parser'
```
