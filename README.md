# sreader: A TUI Atom and RSS Feed Reader

[![Build Status](https://github.com/boneill02/sreader/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/boneill02/sreader/actions/workflows/go.yml).
[![Dependabot Active](https://img.shields.io/badge/dependabot-active-brightgreen?style=flat-square&logo=dependabot)](https://github.com/boneill02/sreader/security/dependabot)

## Usage

1. `go install` (ensure `$HOME/.local/share/go/bin` is in your `$PATH`)
1. Add feed URLs to `~/.config/sreader/urls`
1. Create config file at `~/.config/sreader/config.toml` (optional)
1. Run `sreader sync`
1. Run `sreader`

## Features

- [X] Clean, intuitive TUI interface
- [X] Open entries in browser or media player
- [X] Vim key bindings
- [X] XDG Base Directory Specification compliant

## Keybindings

sreader uses Vim-like keybindings by default.

- `h`: Go back
- `j`: Select next item
- `k`: Select previous item
- `l`: Open selected item
- `o`: Open selected list entry in web browser
- `v`: Open selected list entry in video player
- `r`: Refresh feeds
- `q`: Quit

## Configuration

`sreader` can load color schemes through a config file located at
`~/.config/sreader/config.toml`. See [config_example.toml](config_example.toml)
for an example.

If `Player` or `Browser` are not set in your `config.toml`, `sreader` will
use your `$PLAYER` and `$BROWSER` environment variables as fallbacks. If those
are also not set, `sreader` will default to `mpv` and `firefox` respectively.

The following settings are supported (colors are represented by hex strings):

- `FG`: Primary text color for non-selected list items (and entry contents)
- `BG`: Primary background color for non-selected list items (and entry contents)
- `TitleFG`: Non-selected list entry title foreground color
- `TitleBG`: Non-selected list entry title background color
- `SelectedTitleFG`: Selected list entry title foreground color
- `SelectedTitleBG`: Selected list entry title background color
- `DescFG`: Non-selected list entry description foreground color
- `DescBG`: Non-selected list entry description background color
- `SelectedDescFG`: Selected list entry description foreground color
- `SelectedDescBG`: Selected list entry description background color
- `Player`: Path to default video player
- `Browser`: Path to default browser

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
