# cf-healthy-plugin

## Install

Check the releases page for the [download page](https://github.com/laidbackware/cf-healthy-plugin/releases/). Linux only for now.


## Usage

You must have an active CF CLI session with an account with read privileges across the foundation.

Run `cf health-report` without options to generate `report.xlsx` spreadsheet in the current directory.

```
NAME:
   health-report - Find singleton apps and export them

USAGE:
   cf health-report [OPTIONS]

OPTIONS:
   --format, -f      The format of the output file. (json, xlsx).
   --output, -o      The output file, with or without path.
```

Note: the long interval report only applies to TAS 5.0 onwards.

## Install from Source

```sh
make build && cf install-plugin -f ./bin/healthy-plugin
```

## Uninstall

```sh
cf uninstall-plugin HealthyPlugin
```

## Test plugin

```sh
cf uninstall-plugin HealthyPlugin || True && make build && cf install-plugin -f ./bin/healthy-plugin
```