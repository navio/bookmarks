# Commands

Global flags:

```text
--version        print version
--store <path>   use an alternate bookmarks store
-h, --help       show help
```

## `bm add`

Add a bookmark.

```sh
bm add [name] [path] [--tags a,b,c] [-f|--force]
```

Examples:

```sh
bm add proj . --tags work,Go
bm add
bm add proj .. -f
```

## `bm ls`

List bookmarks (TSV by default).

```sh
bm ls [--json] [--tag x]
```

Examples:

```sh
bm ls
bm ls --json
bm ls --tag work
```

## `bm tags`

List all tags in the store and the number of bookmarks using each tag.

```sh
bm tags [--json]
```

Examples:

```sh
bm tags
bm tags --json
```

## `bm find`

Interactive picker (list). Prints the selected path to stdout.

```sh
bm find [--tag x] [--tags a,b,c]
```

Keys: `enter` print path, `c` copy path, `/` filter, `q` quit.

## `bm table`

Interactive picker (table). Prints the selected path to stdout.

```sh
bm table [--tag x] [--tags a,b,c]
```

Keys: `enter` print path, `c` copy path, `q` quit.

## `bm path`

Print the stored path for a bookmark name.

```sh
bm path <name>
```

Example:

```sh
cd "$(bm path proj)"
```

## `bm go`

Print a shell-safe `cd` command for a bookmark name.

```sh
bm go <name>
```

Example:

```sh
eval "$(bm go proj)"
```

## `bm update`

Rename and/or retag an existing bookmark.

```sh
bm update <name> [--name <new>] [--tags a,b,c]
```

Examples:

```sh
bm update proj --tags work,go,tools
bm update proj --name proj2
```

## `bm rm`

Remove a bookmark.

```sh
bm rm <name> [-f|--force]
```

## `bm init`

Print shell integration that lets your current shell session run `bm go <name>` as a direct directory change.

```sh
bm init [bash|zsh|fish]
```

Examples:

```sh
# zsh/bash
eval "$(bm init zsh)"

# now this changes directory directly
bm go proj

# fish
bm init fish | source

# then use helpers
bmcd
bmcd --tag work
bmgo proj
```

## `bm shell init`

Compatibility alias for `bm init`.

```sh
bm shell init [bash|zsh|fish]
```

## `bm help`

Show help.

```sh
bm help
```
