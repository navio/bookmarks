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
  - title: "Shell jump helpers"
    details: "Run `bm init` to enable direct `bm go <name>` jumps plus `bmcd` and `bmgo` helpers in your shell."
  - title: "Simple store"
    details: "Bookmarks live in a TSV file under your XDG config directory. Human-editable and merge-friendly."
---

## Quickstart

```sh
brew tap navio/tap
brew install navio/tap/bm

bm add proj . --tags work,go
eval "$(bm go proj)"

eval "$(bm init zsh)"
bm go proj
bmgo proj
```
