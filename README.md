# sreader: A TUI Atom and RSS Feed Reader

[![Build Status](https://github.com/boneill02/sreader/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/boneill02/sreader/actions/workflows/go.yml).

## Usage

1. `go install sreader.go`
1. Add feed URLs to `~/.config/sreader/urls`
1. Run `sreader sync`
1. Run `sreader`

## Features

- [X] Clean, intuitive TUI interface
- [X] Open entries in browser or media player
- [X] Vim key bindings
- [X] XDG Base Directory Specification compliant

## Keybindings

sreader uses Vim-like keybindings by default.

* `h`: Go back
* `j`: Select next item
* `k`: Select previous item
* `l`: Open selected item
* `o`: Open selected list entry in `$BROWSER`
* `v`: Open selected list entry in `$PLAYER` (or [mpv](https://mpv.io/) if env
  variable is empty)
* `r`: Refresh feeds
* `q`: Quit

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
