# Store & Format

## Default location

Bookmarks are stored in a TSV file at:

```text
${XDG_CONFIG_HOME:-~/.config}/bm/bookmarks.tsv
```

Override the store file for any command:

```sh
bm --store /tmp/bm.tsv add tmp .
bm --store /tmp/bm.tsv ls
```

## Data format

One bookmark per line:

```text
name\tpath\ttags\tcreated_at
```

- `tags` is a comma-separated list (normalized to lowercase and deduped)
- `created_at` is RFC3339
- blank lines and lines starting with `#` are ignored
