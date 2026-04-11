# Usage Patterns

## Jump to a bookmark

```sh
cd "$(bm path proj)"
eval "$(bm go proj)"
```

## Jump by name with helper

```sh
eval "$(bm init zsh)"   # or bash
bm go proj
bmgo proj
```

## Pick interactively

```sh
cd "$(bm find)"
cd "$(bm table)"
```

## Pick and jump in one step

```sh
eval "$(bm init zsh)"   # or bash
bmcd
bmcd --tag work
```

Fish:

```sh
bm init fish | source
bmcd
```

## Filter by tag

```sh
bm ls --tag work
cd "$(bm find --tag work)"
cd "$(bm table --tags work,go)"
```

## Discover tags

```sh
bm tags
```

## Copy paths from the TUI

In `bm find` / `bm table`, press `c` to copy the selected path to your clipboard.
