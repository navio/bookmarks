# Bookmarks CLI (bm) — Plan

## TL;DR recommendation
- Build in Go for a single static binary and easy Homebrew distribution.
- Use optional fzf for fuzzy selection; provide a non-fzf fallback prompt.
- Provide shell integration via `bm init` that defines a shell function to `cd` into bookmark paths.

## Detailed CLI spec

### Global
- Binary: `bm`
- Default store path: `~/.config/bm/bookmarks.tsv`
- Output: human-friendly by default, with `--json` for machine use where relevant.

### `bm add`
- Purpose: Add or update a bookmark.
- Flags:
  - `-n, --name <name>`: Bookmark name (required unless provided as arg).
  - `-p, --path <path>`: Path to bookmark (default: `.`).
  - `-t, --tags <csv>`: Comma-separated tags.
  - `-f, --force`: Overwrite if name already exists.
- Examples:
  - `bm add -n work -p ~/Projects/work -t client,active`
  - `bm add notes` (adds current dir with name `notes`)

### `bm ls`
- Purpose: List bookmarks.
- Flags:
  - `--tag <tag>`: Filter by a tag (repeatable).
  - `--name <glob>`: Filter by name glob.
  - `--json`: JSON output.
  - `--sort <name|path|created>`: Sort order.
- Examples:
  - `bm ls`
  - `bm ls --tag client --sort name`

### `bm path`
- Purpose: Print a bookmark path (for scripting or `cd`).
- Flags:
  - `--name <name>`: Bookmark name (required unless arg provided).
- Examples:
  - `bm path work`

### `bm pick`
- Purpose: Interactively pick a bookmark and print its path.
- Flags:
  - `--tag <tag>`: Filter before pick.
  - `--name <glob>`: Filter before pick.
  - `--no-fzf`: Disable fzf even if available.
- Examples:
  - `bm pick`
  - `cd "$(bm pick)"`

### `bm rm`
- Purpose: Remove a bookmark by name.
- Flags:
  - `-f, --force`: Skip confirmation.
- Examples:
  - `bm rm work`

### `bm prune`
- Purpose: Remove entries whose paths no longer exist.
- Flags:
  - `-f, --force`: Skip confirmation.
- Examples:
  - `bm prune`

### `bm tags`
- Purpose: List all tags and counts.
- Flags:
  - `--json`: JSON output.
- Examples:
  - `bm tags`

### `bm init`
- Purpose: Print shell integration snippet.
- Behavior: prints a function like `bmcd()` that runs `cd "$(bm path/pick ...)"`.
- Examples:
  - `bm init` (user evals output in shell rc)

### `bm completion`
- Purpose: Generate shell completion scripts.
- Flags:
  - `--shell <bash|zsh|fish>`
- Examples:
  - `bm completion --shell zsh`

## Storage format recommendation
- Preferred: TSV at `~/.config/bm/bookmarks.tsv`.
- Rationale:
  - TSV is human-editable, merge-friendly, and trivial to parse.
  - Each line: `name<TAB>path<TAB>tags<TAB>created_at`.
  - Tags are comma-separated. `created_at` is RFC3339.
- Optional future: JSON export/import for richer metadata.

## Fuzzy selection approach
- Prefer fzf when installed: pipe `bm ls` (plain text) into `fzf` and return selected entry.
- Fallback behavior without fzf:
  - Show a numbered list and prompt for a selection index.
  - Support simple substring filter before list is shown.

## Homebrew distribution plan
- Use GitHub Releases + GoReleaser to build and publish macOS/Linux binaries.
- Create separate tap repo (e.g., `homebrew-bm`) with formula managed by GoReleaser.
- Formula includes:
  - `bin` install.
  - Shell completions (bash/zsh/fish) installed to standard locations.
  - `test do` block verifying `bm --version` and `bm add/ls` in a temp dir.
- Automate release via GitHub Actions on version tags.

## Milestones (PR-sized)
1) Prototype
   - Implement TSV store + `add`, `ls`, `path`, `rm`.
   - Basic error handling and config path resolution.
2) Interactive use
   - Add `pick` with fzf integration + fallback.
   - Add `prune` and `tags`.
3) Shell integration
   - Add `init` output and `completion` command.
4) Release prep
   - Add versioning, CI, and GoReleaser config.
   - Create Homebrew tap + formula + tests.
5) v0.1 release
   - Docs polish, usage examples, and first release tag.

## Open questions
- Should `bm add` allow duplicate names with automatic suffixing, or always require `--force`?
- Should tags be case-sensitive or normalized to lowercase?
- Should `bm ls` default to a compact table or full paths only?
- Do we want `bm edit` (open store in $EDITOR) or `bm mv` (rename) in v0.1?
- Should `bm pick` default to `path` output or optionally `name`?
