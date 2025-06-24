# sreader: A Simple TUI Atom and RSS Feed Reader

![index](https://oneill.sh/img/sreader-index.png)
![feedview](https://oneill.sh/img/sreader-feedview.png)
![entryview](https://oneill.sh/img/sreader-entryview.png)

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

## Bugs

Submit an issue. Email me a patch or submit a PR if you've fixed it.

## License

Copyright (C) 2020-2021 Ben O'Neill <ben@oneill.sh>. License: MIT.
See LICENSE for more details.
