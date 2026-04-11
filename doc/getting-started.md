# Getting Started

`bm` is a tiny CLI that saves and recalls named bookmarks to directories.

## Install

### Homebrew (recommended)

```sh
brew tap navio/tap
brew install navio/tap/bm

bm --version

# enable shell integration (recommended)
echo 'eval "$(bm init zsh)"' >> ~/.zshrc
source ~/.zshrc
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
eval "$(bm go proj)"
```

## Interactive picker

```sh
# list picker (supports '/' filtering)
bm find

# table picker
bm table

# without shell integration, eval the command output
eval "$(bm find)"
```

## Jump directly from your shell

Install shell integration once per shell session:

```sh
# zsh/bash
eval "$(bm init zsh)"

# fish
bm init fish | source
```

Then use:

```sh
bm go proj
bmcd
bmcd --tag work
bmgo proj
```
