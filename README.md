# bookmarks (`bm`)

A tiny CLI to save and recall “bookmarks” to directories.

## Install

### Homebrew (recommended)

```sh
brew tap navio/tap
brew install navio/tap/bm

bm --version
```

### From source

```sh
git clone https://github.com/navio/bookmarks
cd bookmarks

go run ./cmd/bm --version
```

## Usage

```sh
bm help

# add a bookmark
# - if name is omitted, it defaults to the current directory name
# - path defaults to the current directory
bm add proj . --tags work,Go
bm add

# overwrite an existing bookmark (update path; tags only change if provided)
bm add proj .. -f

# list bookmarks
bm ls

# list as JSON
bm ls --json

# filter by tag
bm ls --tag work

# interactive picker (list)
# enter: print path  •  c: copy path  •  /: filter  •  q: quit
bm find
bm find --tags work,go

# interactive table
# enter: print path  •  c: copy path  •  q: quit
bm table
bm table --tag work

# print the path for a bookmark
bm path proj

# print the path for a bookmark (use it like this)
cd "$(bm path proj)"

# update tags and/or rename a bookmark
bm update proj --tags work,go,tools
bm update proj --name proj2

# remove a bookmark
bm rm proj2
```

## Store location

By default, bookmarks are stored in a TSV file at:

- `${XDG_CONFIG_HOME:-~/.config}/bm/bookmarks.tsv`

You can override the store file for any command (useful for testing):

```sh
bm --store /tmp/bm.tsv add tmp .
bm --store /tmp/bm.tsv ls
```

## Data format

The store file is TSV with one entry per line:

```
name\tpath\ttags\tcreated_at
```

- `tags` is a comma-separated list (normalized to lowercase and deduped)
- `created_at` is RFC3339
- blank lines and lines starting with `#` are ignored
