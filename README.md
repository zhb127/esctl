# esctl

A command line tool for easy management of ElasticSearch cluster.

## Install

### MacOS

```bash
brew tap zhb127/esctl

brew install esctl
```

## Usage

```bash
esctl -h
esctl controls the ElasticSearch cluster manager

Usage:
  esctl [command]

Available Commands:
  help        Help about any command
  index       Manage indices
  migrate     Manage migrations
  version     Show the version information

Flags:
      --cluster string   The name of the config cluster to use
  -c, --config string    config file (default is $HOME/.esctl/config)
      --context string   The name of the config context to use
  -h, --help             help for esctl
  -t, --toggle           Help message for toggle
      --user string      The name of the config user to use

Use "esctl [command] --help" for more information about a command.
```