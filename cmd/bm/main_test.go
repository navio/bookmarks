package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/navio/bookmarks/internal/bookmarks"
)

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = old
	}()

	fnErr := fn()
	_ = w.Close()

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String(), fnErr
}

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

func TestCmdTags_EmptyStore(t *testing.T) {
	root := t.TempDir()
	storePath := filepath.Join(root, "bm.tsv")

	out, err := captureStdout(t, func() error {
		return cmdTags(storePath, nil)
	})
	if err != nil {
		t.Fatalf("cmdTags() error = %v", err)
	}
	if out != "" {
		t.Fatalf("stdout=%q, want empty", out)
	}
}

func TestCmdTags_CountsAndSort(t *testing.T) {
	root := t.TempDir()
	storePath := filepath.Join(root, "bm.tsv")

	entries := []bookmarks.Bookmark{
		{Name: "a", Path: "/tmp/a", Tags: []string{"work", "go"}, CreatedAt: time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)},
		{Name: "b", Path: "/tmp/b", Tags: []string{"go"}, CreatedAt: time.Date(2026, 2, 15, 12, 1, 0, 0, time.UTC)},
		{Name: "c", Path: "/tmp/c", Tags: nil, CreatedAt: time.Date(2026, 2, 15, 12, 2, 0, 0, time.UTC)},
	}
	if err := bookmarks.Save(storePath, entries); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	out, err := captureStdout(t, func() error {
		return cmdTags(storePath, []string{})
	})
	if err != nil {
		t.Fatalf("cmdTags() error = %v", err)
	}
	want := "go\t2\nwork\t1\n"
	if out != want {
		t.Fatalf("stdout=\n%q\nwant=\n%q", out, want)
	}
}

func TestCmdTags_JSON(t *testing.T) {
	root := t.TempDir()
	storePath := filepath.Join(root, "bm.tsv")

	entries := []bookmarks.Bookmark{
		{Name: "a", Path: "/tmp/a", Tags: []string{"work", "go"}, CreatedAt: time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)},
		{Name: "b", Path: "/tmp/b", Tags: []string{"go"}, CreatedAt: time.Date(2026, 2, 15, 12, 1, 0, 0, time.UTC)},
	}
	if err := bookmarks.Save(storePath, entries); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	out, err := captureStdout(t, func() error {
		return cmdTags(storePath, []string{"--json"})
	})
	if err != nil {
		t.Fatalf("cmdTags(--json) error = %v", err)
	}

	var got []struct {
		Tag   string `json:"tag"`
		Count int    `json:"count"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("json unmarshal error = %v\nstdout=%q", err, out)
	}
	want := []struct {
		Tag   string `json:"tag"`
		Count int    `json:"count"`
	}{
		{Tag: "go", Count: 2},
		{Tag: "work", Count: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("json=%#v, want %#v", got, want)
	}
}

func TestCmdInit_Bash(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return cmdInit([]string{"bash"})
	})
	if err != nil {
		t.Fatalf("cmdInit(bash) error = %v", err)
	}
	if out == "" {
		t.Fatalf("expected shell script output")
	}
	if !strings.Contains(out, "bm()") {
		t.Fatalf("expected bm wrapper function in output, got %q", out)
	}
	if !strings.Contains(out, "command bm go") {
		t.Fatalf("expected bm go interception in output, got %q", out)
	}
	if !strings.Contains(out, "bmcd()") {
		t.Fatalf("expected bmcd function in output, got %q", out)
	}
	if !strings.Contains(out, "bmgo()") {
		t.Fatalf("expected bmgo function in output, got %q", out)
	}
}

func TestCmdInit_AutodetectShell(t *testing.T) {
	t.Setenv("SHELL", "/bin/zsh")
	out, err := captureStdout(t, func() error {
		return cmdInit([]string{})
	})
	if err != nil {
		t.Fatalf("cmdInit() error = %v", err)
	}
	if !strings.Contains(out, "bmcd()") {
		t.Fatalf("expected sh-compatible output, got %q", out)
	}
}

func TestCmdInit_UnsupportedShell(t *testing.T) {
	err := cmdInit([]string{"pwsh"})
	if err == nil {
		t.Fatalf("expected unsupported shell error")
	}
}

func TestCmdShellInit_CompatAlias(t *testing.T) {
	out, err := captureStdout(t, func() error {
		return cmdShell([]string{"init", "bash"})
	})
	if err != nil {
		t.Fatalf("cmdShell(init bash) error = %v", err)
	}
	if !strings.Contains(out, "bm()") {
		t.Fatalf("expected bm wrapper in compatibility output, got %q", out)
	}
}

func TestCmdGo_PrintsCdCommand(t *testing.T) {
	root := t.TempDir()
	storePath := filepath.Join(root, "bm.tsv")

	entries := []bookmarks.Bookmark{
		{Name: "proj", Path: "/tmp/my proj", CreatedAt: time.Date(2026, 4, 10, 12, 0, 0, 0, time.UTC)},
	}
	if err := bookmarks.Save(storePath, entries); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	out, err := captureStdout(t, func() error {
		return cmdGo(storePath, []string{"proj"})
	})
	if err != nil {
		t.Fatalf("cmdGo() error = %v", err)
	}
	want := "cd -- '/tmp/my proj'\n"
	if out != want {
		t.Fatalf("stdout=%q, want %q", out, want)
	}
}

func TestCmdGo_NotFound(t *testing.T) {
	root := t.TempDir()
	storePath := filepath.Join(root, "bm.tsv")

	err := cmdGo(storePath, []string{"missing"})
	if err == nil {
		t.Fatalf("expected bookmark not found error")
	}
}
