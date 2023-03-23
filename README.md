[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/jadiunr/check-smart-status)
![Go Test](https://github.com/jadiunr/check-smart-status/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/jadiunr/check-smart-status/workflows/goreleaser/badge.svg)

# Sensu S.M.A.R.T. status check

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu S.M.A.R.T. status check is a [Sensu Check][6] that provide S.M.A.R.T. statistics for all storage devices.

The user running Check must have permission to access the storage device. (e.g. Add `sudo` as prefix to the command)

The only supported interfaces are SATA and NVMe.
SCSI and other hardware RAID are not supported.

## Usage examples

```
S.M.A.R.T. status check for Sensu

Usage:
  check-smart-status [flags]
  check-smart-status [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -h, --help   help for check-smart-status

Use "check-smart-status [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add jadiunr/check-smart-status
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/jadiunr/check-smart-status].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-smart-status
  namespace: default
spec:
  command: sudo check-smart-status
  subscriptions:
  - system
  runtime_assets:
  - jadiunr/check-smart-status
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the check-smart-status repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/check-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/check-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[7]: https://github.com/sensu-community/check-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
