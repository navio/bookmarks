# Usage Patterns

## Jump to a bookmark

```sh
cd "$(bm path proj)"
```

## Pick interactively

```sh
cd "$(bm find)"
cd "$(bm table)"
```

## Filter by tag

```sh
bm ls --tag work
cd "$(bm find --tag work)"
cd "$(bm table --tags work,go)"
```

## Copy paths from the TUI

In `bm find` / `bm table`, press `c` to copy the selected path to your clipboard.
