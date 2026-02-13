# Getting Started

`bm` is a tiny CLI that saves and recalls named bookmarks to directories.

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

## First bookmarks

```sh
# name defaults to current directory name; path defaults to '.'
bm add

bm add proj . --tags work,go
bm ls

cd "$(bm path proj)"
```

## Interactive picker

```sh
# list picker (supports '/' filtering)
cd "$(bm find)"

# table picker
cd "$(bm table)"
```
