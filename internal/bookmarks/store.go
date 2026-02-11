package bookmarks

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Bookmark represents a single entry in the TSV store.
type Bookmark struct {
	Name      string
	Path      string
	Tags      []string
	CreatedAt time.Time
}

// DefaultPath returns the default TSV storage path.
func DefaultPath() (string, error) {
	if xdg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); xdg != "" {
		return filepath.Join(xdg, "bm", "bookmarks.tsv"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".config", "bm", "bookmarks.tsv"), nil
}

// Load reads bookmarks from a TSV file. Missing files return an empty slice.
func Load(path string) ([]Bookmark, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Bookmark{}, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	entries := []Bookmark{}
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			return nil, fmt.Errorf("line %d: expected 4 fields", lineNum)
		}
		createdAt, err := time.Parse(time.RFC3339, parts[3])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse created_at: %w", lineNum, err)
		}
		entry := Bookmark{
			Name:      parts[0],
			Path:      parts[1],
			Tags:      normalizeTags(parts[2]),
			CreatedAt: createdAt,
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// Save writes bookmarks to a TSV file atomically.
func Save(path string, entries []Bookmark) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "bookmarks-*.tsv")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	writer := bufio.NewWriter(tmp)
	for _, entry := range entries {
		line := fmt.Sprintf("%s\t%s\t%s\t%s\n",
			entry.Name,
			entry.Path,
			tagsToString(entry.Tags),
			entry.CreatedAt.Format(time.RFC3339),
		)
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// NormalizeTags converts a comma-separated tag string into normalized tags.
func NormalizeTags(input string) []string {
	return normalizeTags(input)
}

func normalizeTags(input string) []string {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	seen := map[string]struct{}{}
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := strings.ToLower(strings.TrimSpace(part))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func tagsToString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return strings.Join(tags, ",")
}

// ResolvePath returns an absolute, cleaned path.
func ResolvePath(input string, cwd string) (string, error) {
	path := input
	if strings.TrimSpace(path) == "" {
		path = cwd
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

// ContainsTag reports whether tags contain the target tag.
func ContainsTag(tags []string, target string) bool {
	needle := strings.ToLower(strings.TrimSpace(target))
	if needle == "" {
		return false
	}
	for _, tag := range tags {
		if strings.ToLower(tag) == needle {
			return true
		}
	}
	return false
}
