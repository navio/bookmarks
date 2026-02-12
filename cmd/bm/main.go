package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/navio/bookmarks/internal/bookmarks"
)

const version = "0.2.1"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	opts, rest, err := parseGlobalArgs(args)
	if err != nil {
		return err
	}
	if opts.version {
		fmt.Printf("bm %s\n", version)
		return nil
	}
	if opts.help || len(rest) == 0 {
		fmt.Println(usage())
		return nil
	}

	storePath := opts.storePath
	if storePath == "" {
		p, err := bookmarks.DefaultPath()
		if err != nil {
			return err
		}
		storePath = p
	}
	storePath, err = resolveStorePath(storePath)
	if err != nil {
		return err
	}

	switch rest[0] {
	case "add":
		return cmdAdd(storePath, rest[1:])
	case "ls":
		return cmdList(storePath, rest[1:])
	case "find":
		return cmdFind(storePath, rest[1:])
	case "table":
		return cmdTable(storePath, rest[1:])
	case "path":
		return cmdPath(storePath, rest[1:])
	case "update":
		return cmdUpdate(storePath, rest[1:])
	case "rm":
		return cmdRemove(storePath, rest[1:])
	case "help":
		fmt.Println(usage())
		return nil
	default:
		return fmt.Errorf("unknown command: %s\n\n%s", rest[0], usage())
	}
}

func cmdAdd(storePath string, args []string) error {
	positionals, err := parseArgs(args, map[string]bool{"--tags": true, "-f": false, "--force": false})
	if err != nil {
		return err
	}
	var (
		tagsInput string
		hasTags   bool
	)
	if value, ok := positionals.flags["--tags"]; ok {
		tagsInput = value
		hasTags = true
	}
	_, forceShort := positionals.flags["-f"]
	_, forceLong := positionals.flags["--force"]
	force := forceShort || forceLong

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if len(positionals.args) > 2 {
		return fmt.Errorf("usage: bm add [name] [path] [--tags a,b,c] [-f|--force]")
	}

	name := ""
	pathInput := ""
	if len(positionals.args) >= 1 {
		name = strings.TrimSpace(positionals.args[0])
	}
	if len(positionals.args) == 2 {
		pathInput = positionals.args[1]
	}
	if name == "" {
		name = filepath.Base(cwd)
	}
	if name == "" || name == "." || name == string(filepath.Separator) {
		return errors.New("name cannot be empty")
	}
	if strings.Contains(name, "\t") || strings.Contains(name, "\n") {
		return errors.New("name cannot contain tabs or newlines")
	}

	resolvedPath, err := bookmarks.ResolvePath(pathInput, cwd)
	if err != nil {
		return err
	}

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}

	for i := range entries {
		if entries[i].Name != name {
			continue
		}
		if !force {
			return fmt.Errorf("bookmark already exists: %s", name)
		}
		entries[i].Path = resolvedPath
		if hasTags {
			entries[i].Tags = bookmarks.NormalizeTags(tagsInput)
		}
		return bookmarks.Save(storePath, entries)
	}

	entry := bookmarks.Bookmark{
		Name:      name,
		Path:      resolvedPath,
		Tags:      nil,
		CreatedAt: time.Now().UTC(),
	}
	if hasTags {
		entry.Tags = bookmarks.NormalizeTags(tagsInput)
	}
	entries = append(entries, entry)

	if err := bookmarks.Save(storePath, entries); err != nil {
		return err
	}
	return nil
}

func cmdList(storePath string, args []string) error {
	positionals, err := parseArgs(args, map[string]bool{"--json": false, "--tag": true})
	if err != nil {
		return err
	}
	if len(positionals.args) != 0 {
		return errors.New("usage: bm ls [--json] [--tag x]")
	}
	_, jsonOutput := positionals.flags["--json"]
	tagFilter := positionals.flags["--tag"]

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}

	filtered := make([]bookmarks.Bookmark, 0, len(entries))
	for _, entry := range entries {
		if tagFilter != "" && !bookmarks.ContainsTag(entry.Tags, tagFilter) {
			continue
		}
		filtered = append(filtered, entry)
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name < filtered[j].Name
	})

	if jsonOutput {
		payload := make([]map[string]any, 0, len(filtered))
		for _, entry := range filtered {
			payload = append(payload, map[string]any{
				"name":       entry.Name,
				"path":       entry.Path,
				"tags":       entry.Tags,
				"created_at": entry.CreatedAt.Format(time.RFC3339),
			})
		}
		encoded, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(encoded))
		return nil
	}

	for _, entry := range filtered {
		fmt.Printf("%s\t%s\t%s\t%s\n",
			entry.Name,
			entry.Path,
			strings.Join(entry.Tags, ","),
			entry.CreatedAt.Format(time.RFC3339),
		)
	}
	return nil
}

func cmdFind(storePath string, args []string) error {
	positionals, err := parseArgs(args, map[string]bool{"--tag": true, "--tags": true})
	if err != nil {
		return err
	}
	if len(positionals.args) != 0 {
		return errors.New("usage: bm find [--tag x] [--tags a,b,c]")
	}

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}
	tags := parseTagFilters(positionals.flags)
	entries = filterByAnyTag(entries, tags)

	selected, err := runFindTUI(entries, "bm find")
	if err != nil {
		return err
	}
	if strings.TrimSpace(selected) != "" {
		fmt.Println(selected)
	}
	return nil
}

func cmdTable(storePath string, args []string) error {
	positionals, err := parseArgs(args, map[string]bool{"--tag": true, "--tags": true})
	if err != nil {
		return err
	}
	if len(positionals.args) != 0 {
		return errors.New("usage: bm table [--tag x] [--tags a,b,c]")
	}

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}
	tags := parseTagFilters(positionals.flags)
	entries = filterByAnyTag(entries, tags)

	selected, err := runTableTUI(entries, "bm table")
	if err != nil {
		return err
	}
	if strings.TrimSpace(selected) != "" {
		fmt.Println(selected)
	}
	return nil
}

func parseTagFilters(flags map[string]string) []string {
	out := []string{}
	if v, ok := flags["--tag"]; ok {
		out = append(out, v)
	}
	if v, ok := flags["--tags"]; ok {
		// supports comma-separated list
		out = append(out, v)
	}
	if len(out) == 0 {
		return nil
	}

	// Normalize and dedupe.
	merged := strings.Join(out, ",")
	return bookmarks.NormalizeTags(merged)
}

func filterByAnyTag(entries []bookmarks.Bookmark, tags []string) []bookmarks.Bookmark {
	if len(tags) == 0 {
		return entries
	}
	filtered := make([]bookmarks.Bookmark, 0, len(entries))
	for _, e := range entries {
		for _, t := range tags {
			if bookmarks.ContainsTag(e.Tags, t) {
				filtered = append(filtered, e)
				break
			}
		}
	}
	return filtered
}

func cmdPath(storePath string, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: bm path <name>")
	}
	name := args[0]

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name == name {
			fmt.Println(entry.Path)
			return nil
		}
	}
	return fmt.Errorf("bookmark not found: %s", name)
}

func cmdUpdate(storePath string, args []string) error {
	positionals, err := parseArgs(args, map[string]bool{"--name": true, "--tags": true})
	if err != nil {
		return err
	}
	if len(positionals.args) != 1 {
		return errors.New("usage: bm update <name> [--name <new>] [--tags a,b,c]")
	}
	oldName := strings.TrimSpace(positionals.args[0])
	if oldName == "" {
		return errors.New("name cannot be empty")
	}

	newNameRaw, hasNewName := positionals.flags["--name"]
	tagsRaw, hasTags := positionals.flags["--tags"]

	if !hasNewName && !hasTags {
		return errors.New("nothing to update: provide --name and/or --tags")
	}

	newName := strings.TrimSpace(newNameRaw)
	if hasNewName {
		if newName == "" {
			return errors.New("new name cannot be empty")
		}
		if strings.Contains(newName, "\t") || strings.Contains(newName, "\n") {
			return errors.New("new name cannot contain tabs or newlines")
		}
	}

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}

	if hasNewName {
		for _, entry := range entries {
			if entry.Name == newName && entry.Name != oldName {
				return fmt.Errorf("bookmark already exists: %s", newName)
			}
		}
	}

	found := false
	for i := range entries {
		if entries[i].Name != oldName {
			continue
		}
		found = true
		if hasNewName {
			entries[i].Name = newName
		}
		if hasTags {
			entries[i].Tags = bookmarks.NormalizeTags(tagsRaw)
		}
	}

	if !found {
		return fmt.Errorf("bookmark not found: %s", oldName)
	}
	return bookmarks.Save(storePath, entries)
}

func cmdRemove(storePath string, args []string) error {
	positionals, err := parseArgs(args, map[string]bool{"-f": false, "--force": false})
	if err != nil {
		return err
	}
	if len(positionals.args) != 1 {
		return errors.New("usage: bm rm <name> [-f|--force]")
	}
	name := positionals.args[0]
	_, forceShort := positionals.flags["-f"]
	_, forceLong := positionals.flags["--force"]
	force := forceShort || forceLong

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		return err
	}

	result := make([]bookmarks.Bookmark, 0, len(entries))
	removed := false
	for _, entry := range entries {
		if entry.Name == name {
			removed = true
			continue
		}
		result = append(result, entry)
	}

	if !removed {
		if force {
			return nil
		}
		return fmt.Errorf("bookmark not found: %s", name)
	}

	return bookmarks.Save(storePath, result)
}

type parsedArgs struct {
	args  []string
	flags map[string]string
}

// parseArgs extracts flags (with optional values) and positional args.
// allowed maps flag name to whether it expects a value.
func parseArgs(args []string, allowed map[string]bool) (parsedArgs, error) {
	parsed := parsedArgs{
		flags: map[string]string{},
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			name, value, hasValue := strings.Cut(arg, "=")
			expectsValue, ok := allowed[name]
			if !ok {
				return parsed, fmt.Errorf("unknown flag: %s", name)
			}
			if expectsValue {
				if hasValue {
					parsed.flags[name] = value
					continue
				}
				if i+1 >= len(args) {
					return parsed, fmt.Errorf("flag %s requires a value", name)
				}
				i++
				parsed.flags[name] = args[i]
				continue
			}
			if hasValue {
				return parsed, fmt.Errorf("flag %s does not take a value", name)
			}
			parsed.flags[name] = ""
			continue
		}
		if strings.HasPrefix(arg, "-") {
			expectsValue, ok := allowed[arg]
			if !ok {
				return parsed, fmt.Errorf("unknown flag: %s", arg)
			}
			if expectsValue {
				if i+1 >= len(args) {
					return parsed, fmt.Errorf("flag %s requires a value", arg)
				}
				i++
				parsed.flags[arg] = args[i]
				continue
			}
			parsed.flags[arg] = ""
			continue
		}
		parsed.args = append(parsed.args, arg)
	}
	return parsed, nil
}

func usage() string {
	return strings.TrimSpace(`usage:
  bm --version
  bm [--store <path>] <command>
  bm add [name] [path] [--tags a,b,c] [-f|--force]
  bm ls [--json] [--tag x]
  bm find [--tag x] [--tags a,b,c]
  bm table [--tag x] [--tags a,b,c]
  bm path <name>
  bm update <name> [--name <new>] [--tags a,b,c]
  bm rm <name> [-f|--force]

global flags:
  --store <path>   override default store path
  -h, --help       show help
  --version        print version`)
}

type globalOpts struct {
	storePath string
	help      bool
	version   bool
}

func parseGlobalArgs(args []string) (globalOpts, []string, error) {
	opts := globalOpts{}
	rest := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--version":
			opts.version = true
		case arg == "-h" || arg == "--help":
			opts.help = true
		case arg == "--store":
			if i+1 >= len(args) {
				return opts, nil, errors.New("flag --store requires a value")
			}
			i++
			opts.storePath = args[i]
		case strings.HasPrefix(arg, "--store="):
			_, v, _ := strings.Cut(arg, "=")
			if strings.TrimSpace(v) == "" {
				return opts, nil, errors.New("flag --store requires a value")
			}
			opts.storePath = v
		default:
			rest = append(rest, arg)
		}
	}
	return opts, rest, nil
}

func resolveStorePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", errors.New("store path cannot be empty")
	}
	if filepath.IsAbs(trimmed) {
		return filepath.Clean(trimmed), nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Clean(filepath.Join(cwd, trimmed)), nil
}
