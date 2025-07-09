# sreader: A TUI Atom and RSS Feed Reader

[![Build Status](https://github.com/bmoneill/sreader/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/bmoneill/sreader/actions/workflows/go.yml).
[![Dependabot Active](https://img.shields.io/badge/dependabot-active-brightgreen?style=flat-square&logo=dependabot)](https://github.com/bmoneill/sreader/security/dependabot)

## Installation

```shell
go install github.com/bmoneill/sreader@latest
$GOPATH/bin/sreader
```

## Usage

```shell
sreader [-c configfile] [-s]
```

- `-c`: Set configuration file
- `-s`: Sync feeds

## Features

- [X] Clean, intuitive TUI interface
- [X] Open entries in browser or media player
- [X] Vim key bindings
- [X] [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/latest/) compliant

## Keybindings

sreader uses Vim-like keybindings by default.

- `h`: Go back
- `j`: Select next item
- `k`: Select previous item
- `l`: Open selected item
- `/`: Filter list items
- `o`: Open selected list entry in web browser
- `v`: Open selected list entry in video player
- `r`: Refresh feeds
- `q`: Quit

## Configuration

sreader can load settings, including feed URLs, colors, keybindings, and paths
through a configuration file located at `~/sreader/config.toml`.
See [config_example.toml](config_example.toml) for an example.

The example configuration file contains all the default values, besides the URL
list. After installing, you can run sreader to generate a configuration file
and then add your feed URLs. Colors must be in hex format.

sreader will also use `$BROWSER` and `$PLAYER` environment variables if not
overridden by your configuration file.

If `$XDG_CONFIG_HOME` is set, sreader will load config files at
`$XDG_CONFIG_HOME/sreader/sreader.toml` by default.

If `$XDG_DATA_HOME` is set (and these paths are not overridden in your
configuration file), sreader will default to the following paths:

- `DBFile`: `$XDG_DATA_HOME/sreader/sreader.db`
- `LogFile`: `$XDG_DATA_HOME/sreader/sreader.log`
- `TmpDir`: `$XDG_DATA_HOME/sreader`

## Screenshots

### Index view

![index](https://oneill.sh/img/sreader-index.png)

### Feed view

![feedview](https://oneill.sh/img/sreader-feedview.png)

### Entry view

![entryview](https://oneill.sh/img/sreader-entryview.png)

## Bugs

If you find a bug, submit an issue, PR, or email me with a description and/or patch.

## License

Copyright (c) 2020-2025 Ben O'Neill <ben@oneill.sh>. This work is released under the
terms of the MIT License. See [LICENSE](LICENSE) for the license terms.
