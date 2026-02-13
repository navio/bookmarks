---
layout: home

hero:
  name: "bm"
  text: "directory bookmarks, instantly"
  tagline: "Save a name for a path. Jump back with a single command."
  actions:
    - theme: brand
      text: Getting Started
      link: /getting-started
    - theme: alt
      text: Commands
      link: /commands

features:
  - title: "Fast path recall"
    details: "Use `bm path <name>` for scripts and `cd \"$(bm path name)\"` in your shell."
  - title: "TUI picker"
    details: "Use `bm find` or `bm table` to fuzzy-pick, print, or copy paths."
  - title: "Simple store"
    details: "Bookmarks live in a TSV file under your XDG config directory. Human-editable and merge-friendly."
---

## Quickstart

```sh
brew tap navio/tap
brew install navio/tap/bm

bm add proj . --tags work,go
cd "$(bm path proj)"
```

## Author

Docs maintained by [Alberto Navarro](https://alberto.pub).
