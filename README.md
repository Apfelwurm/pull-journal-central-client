# pull-journal-central-client
[![Test, Build and publish app release](https://github.com/Apfelwurm/pull-journal-central-client/actions/workflows/test-and-build.yml/badge.svg)](https://github.com/Apfelwurm/pull-journal-central-client/actions/workflows/test-and-build.yml)

This is the client command line tool for [pull-journal-central](https://github.com/Apfelwurm/pull-journal-central).

## Installation

Download the latest release [here](https://github.com/Apfelwurm/pull-journal-central-client/releases/latest).
You can find either a deb package which can be installed by  `dpkg -i pull-journal-central-client_*_amd64.deb` or a tar.gz file that can be unpacked using `tar xvf pull-journal-central-client-*-linux-amd64.tar.gz`


## Command Usage

```
Usage:
  pull-journal-central-client [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  log         Create a log entry
  register    Register a device

Flags:
  -h, --help      help for pull-journal-central-client
  -v, --version   version for pull-journal-central-client

Use "pull-journal-central-client [command] --help" for more information about a command.
```

### register command

```
Register a device

Usage:
  pull-journal-central-client register [flags]

Flags:
      --baseURL string                base url of the pjc installation
  -h, --help                          help for register
      --name string                   Name
      --organisationID string         Organisation ID
      --organisationpassword string   Organisation Password
```


### log command

```
Create a log entry

Usage:
  pull-journal-central-client log [flags]

Flags:
      --baseURL string   base url of the pjc installation
      --class string     class of the Log Entry
  -h, --help             help for log
      --service string   service name
      --source string    source of the log Entry
```
