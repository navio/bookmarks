# bookmarks

CLI directory bookmarks (bm)

## Usage

```sh
bm --version
bm help

# optionally override the store file path
bm --store /tmp/bookmarks.tsv ls

bm add <name> [path] [--tags a,b,c]
bm ls [--json] [--tag x]
bm path <name>
bm rm <name> [-f|--force]
```

Store location:
- Default: `${XDG_CONFIG_HOME:-~/.config}/bm/bookmarks.tsv`
- Override: `bm --store <path> ...`
