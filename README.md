# sreader: A Simple TUI Atom and RSS Feed Reader

## Usage

1. `go install sreader.go`
1. Add feed URLs to `~/.config/sreader/urls`
1. Run `sreader sync`
1. Run `sreader`

## Features

- [X] Simple, clean TUI interface
- [X] Open in browser or media player
- [X] Vim keys
- [X] XDG Base Directory Specification compliant
- [ ] Show when an entry has been read
- [ ] View all entries at once
- [ ] Filter out read entries
- [ ] Regex search through entries
- [ ] "Nickname" feeds
- [ ] Config file for keybindings and default browsers/players
- [ ] Search for feeds
- [ ] Color theming
- [ ] Add feeds through the TUI

## Keybindings

sreader uses Vim-like keybindings by default.

* `l`: Open selected list entry in sreader
* `j`: Select next list entry
* `k`: Select previous list entry
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

Copyright (c) 2019-2025 Ben O'Neill <ben@oneill.sh>. This work is released under the
terms of the MIT License. See [LICENSE](LICENSE) for the license terms.
