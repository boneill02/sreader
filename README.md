# sreader: A TUI Atom and RSS Feed Reader

[![Build Status](https://github.com/bmoneill/sreader/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/bmoneill/sreader/actions/workflows/go.yml).
[![Dependabot Active](https://img.shields.io/badge/dependabot-active-brightgreen?style=flat-square&logo=dependabot)](https://github.com/bmoneill/sreader/security/dependabot)

## Usage

1. `go install` (ensure `$HOME/.local/share/go/bin` is in your `$PATH`)
1. Create config file at `~/.config/sreader/config.toml` and add URLs
1. Run `sreader sync`
1. Run `sreader`

## Features

- [X] Clean, intuitive TUI interface
- [X] Open entries in browser or media player
- [X] Vim key bindings

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

`sreader` can load settings through a config file located at
`~/.config/sreader/config.toml`. See [config_example.toml](config_example.toml)
for an example (everything there except the URL list are the defaults).

A config file with the URL list must be present before running `sreader`.

The following settings are supported (colors are represented by hex strings):

- `URLs` (**REQUIRED**): List of feed URLs

### Paths

- `DBFile`: Path to the database
- `LogFile`: Path to the log file
- `TmpDir`: Where to store temporary files (used during sync)

### Colors

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

### Keys

- `UpKey`: Move up
- `DownKey`: Move down
- `LeftKey`: Move left
- `RightKey`: Move right
- `QuitKey`: Quit
- `SyncKey`: Sync feeds
- `BrowserKey`: Open entry in browser
- `PlayerKey`: Open entry in media player

### External applications

If these are not set in your configuration file but the `$BROWSER` or `$PLAYER`
environment variables are set, those will be used respectively.

- `Browser`: Path to default browser
- `Player`: Path to default video player

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
