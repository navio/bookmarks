package bookmarks

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestDefaultPath_UsesXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-config")
	p, err := DefaultPath()
	if err != nil {
		t.Fatalf("DefaultPath() error = %v", err)
	}
	want := filepath.Join("/tmp/xdg-config", "bm", "bookmarks.tsv")
	if p != want {
		t.Fatalf("DefaultPath() = %q, want %q", p, want)
	}
}

func TestNormalizeTags(t *testing.T) {
	got := NormalizeTags(" Work,go, ,GO,tools ")
	want := []string{"work", "go", "tools"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeTags() = %#v, want %#v", got, want)
	}
}

func TestSaveLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.tsv")
	entries := []Bookmark{
		{
			Name:      "a",
			Path:      "/tmp/a",
			Tags:      []string{"work"},
			CreatedAt: time.Date(2026, 2, 11, 12, 0, 0, 0, time.UTC),
		},
		{
			Name:      "b",
			Path:      "/tmp/b",
			Tags:      nil,
			CreatedAt: time.Date(2026, 2, 11, 12, 1, 0, 0, time.UTC),
		},
	}

	if err := Save(path, entries); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Prepend a comment and a blank line to ensure Load ignores them.
	orig, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	prefixed := append([]byte("# comment\n\n"), orig...)
	if err := os.WriteFile(path, prefixed, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(got) != len(entries) {
		t.Fatalf("Load() len=%d, want %d", len(got), len(entries))
	}
	for i := range entries {
		if got[i].Name != entries[i].Name || got[i].Path != entries[i].Path || !got[i].CreatedAt.Equal(entries[i].CreatedAt) {
			t.Fatalf("Load()[%d] = %#v, want %#v", i, got[i], entries[i])
		}
		if !reflect.DeepEqual(got[i].Tags, entries[i].Tags) {
			t.Fatalf("Load()[%d].Tags = %#v, want %#v", i, got[i].Tags, entries[i].Tags)
		}
	}
}
