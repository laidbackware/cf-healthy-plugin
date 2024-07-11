# cf-healthy-plugin


## Install

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