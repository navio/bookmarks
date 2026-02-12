package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/navio/bookmarks/internal/bookmarks"
)

func TestCmdAdd_NoNameUsesCurrentDirBase(t *testing.T) {
	root := t.TempDir()
	projDir := filepath.Join(root, "myproj")
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(projDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	storePath := filepath.Join(root, "bm.tsv")
	if err := cmdAdd(storePath, []string{}); err != nil {
		t.Fatalf("cmdAdd() error = %v", err)
	}

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("len(entries)=%d, want 1", len(entries))
	}
	if entries[0].Name != "myproj" {
		t.Fatalf("Name=%q, want %q", entries[0].Name, "myproj")
	}
	gotPath, err := filepath.EvalSymlinks(entries[0].Path)
	if err != nil {
		t.Fatalf("EvalSymlinks(got) error = %v", err)
	}
	wantPath, err := filepath.EvalSymlinks(projDir)
	if err != nil {
		t.Fatalf("EvalSymlinks(want) error = %v", err)
	}
	if gotPath != wantPath {
		t.Fatalf("Path=%q, want %q", gotPath, wantPath)
	}
}

func TestCmdAdd_OverwriteRequiresForce(t *testing.T) {
	root := t.TempDir()
	storePath := filepath.Join(root, "bm.tsv")

	// Create initial bookmark.
	oldCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(oldCwd) }()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := cmdAdd(storePath, []string{"proj"}); err != nil {
		t.Fatalf("cmdAdd initial error = %v", err)
	}

	// Move to a new directory and try to overwrite.
	newDir := filepath.Join(root, "new")
	if err := os.MkdirAll(newDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.Chdir(newDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := cmdAdd(storePath, []string{"proj"}); err == nil {
		t.Fatalf("expected error without force")
	}
	if err := cmdAdd(storePath, []string{"proj", "--force"}); err != nil {
		t.Fatalf("expected overwrite with --force, got %v", err)
	}

	entries, err := bookmarks.Load(storePath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("len(entries)=%d, want 1", len(entries))
	}
	gotPath, err := filepath.EvalSymlinks(entries[0].Path)
	if err != nil {
		t.Fatalf("EvalSymlinks(got) error = %v", err)
	}
	wantPath, err := filepath.EvalSymlinks(newDir)
	if err != nil {
		t.Fatalf("EvalSymlinks(want) error = %v", err)
	}
	if gotPath != wantPath {
		t.Fatalf("Path=%q, want %q", gotPath, wantPath)
	}
}
