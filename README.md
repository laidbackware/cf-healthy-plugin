# cf-healthy-plugin

This CLI plugin add 2 commands:
- `sig-check` - performs a rolling restart of an app and checks whether it responds to the SIGTERM signal, or whether the system has to issue a SIGKILL to terminate the app process.
- `health-report` - generates a report against all apps the user has visibility of, showing how health checks are configured.

## Install

Check the releases page for the [releases page](https://github.com/laidbackware/cf-healthy-plugin/releases/). Linux only for now.

Install the plugin:

```sh
cf install-plugin healthy-plugin-linux-amd64-vx.x.x
```

## Usage `sig-check`

You must have an active CF CLI session with an account with the ability to trigger a deployment and read from log cache.

[This demo app](https://github.com/laidbackware/go-sigterm-ignore) can be use to prove how an app shouldn't behave and the expected outputs.

Run `cf sig-check <app-name>` to perform a rolling restart of the app.

```sh
NAME:
   sig-check - Rolling restart app and check that no SIGKILLs were sent

USAGE:
   cf sig-check [OPTIONS] <app-name>

OPTIONS:
   --debug, -d      Enabled debug logging, which will desplay all app logs.
```

## Usage `health-report`

You must have an active CF CLI session with an account with read privileges across the foundation.

Run `cf health-report` without options to generate `report.xlsx` spreadsheet in the current directory.

```sh
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