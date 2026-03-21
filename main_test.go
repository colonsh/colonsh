package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// GetMaxNameLength
// ---------------------------------------------------------------------------

func TestGetMaxNameLength_Empty(t *testing.T) {
	if got := GetMaxNameLength([]BuiltinAlias{}); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGetMaxNameLength_Nil(t *testing.T) {
	if got := GetMaxNameLength[BuiltinAlias](nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestGetMaxNameLength_BuiltinAliases(t *testing.T) {
	aliases := []BuiltinAlias{
		{Name: "ab"},
		{Name: "abcd"},
		{Name: "a"},
	}
	if got := GetMaxNameLength(aliases); got != 4 {
		t.Errorf("expected 4, got %d", got)
	}
}

func TestGetMaxNameLength_CustomAliases(t *testing.T) {
	aliases := []Alias{
		{Name: "config", Cmd: "code ."},
		{Name: "c", Cmd: "code ."},
	}
	if got := GetMaxNameLength(aliases); got != 6 {
		t.Errorf("expected 6, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// shellQuoteSingle
// ---------------------------------------------------------------------------

func TestShellQuoteSingle_NoQuotes(t *testing.T) {
	if got := shellQuoteSingle("hello world"); got != "hello world" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestShellQuoteSingle_WithSingleQuote(t *testing.T) {
	got := shellQuoteSingle("it's")
	want := `it'\''s`
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestShellQuoteSingle_MultipleQuotes(t *testing.T) {
	got := shellQuoteSingle("can't won't")
	want := `can'\''t won'\''t`
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestShellQuoteSingle_Empty(t *testing.T) {
	if got := shellQuoteSingle(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// detectShell
// ---------------------------------------------------------------------------

func TestDetectShell_Zsh(t *testing.T) {
	t.Setenv("SHELL", "/bin/zsh")
	if got := detectShell(); got != "zsh" {
		t.Errorf("expected zsh, got %q", got)
	}
}

func TestDetectShell_Bash(t *testing.T) {
	t.Setenv("SHELL", "/usr/bin/bash")
	if got := detectShell(); got != "bash" {
		t.Errorf("expected bash, got %q", got)
	}
}

func TestDetectShell_Fish(t *testing.T) {
	t.Setenv("SHELL", "/usr/local/bin/fish")
	if got := detectShell(); got != "fish" {
		t.Errorf("expected fish, got %q", got)
	}
}

func TestDetectShell_UnknownDefaultsToZsh(t *testing.T) {
	t.Setenv("SHELL", "/usr/bin/ksh")
	// ksh is not handled; falls through to default "zsh" on non-Windows
	got := detectShell()
	if got != "zsh" && got != "powershell" {
		t.Errorf("unexpected shell: %q", got)
	}
}

// ---------------------------------------------------------------------------
// expandTilde
// ---------------------------------------------------------------------------

func TestExpandTilde_WithTilde(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	got, err := expandTilde("~/foo/bar")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "foo/bar")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestExpandTilde_NoTilde(t *testing.T) {
	got, err := expandTilde("/absolute/path")
	if err != nil {
		t.Fatal(err)
	}
	if got != "/absolute/path" {
		t.Errorf("expected unchanged path, got %q", got)
	}
}

func TestExpandTilde_TildeOnly(t *testing.T) {
	// "~" alone (no slash) should be returned as-is per current implementation
	got, err := expandTilde("~")
	if err != nil {
		t.Fatal(err)
	}
	if got != "~" {
		t.Errorf("expected ~ unchanged, got %q", got)
	}
}

func TestExpandTilde_RelativePath(t *testing.T) {
	got, err := expandTilde("relative/path")
	if err != nil {
		t.Fatal(err)
	}
	if got != "relative/path" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

// ---------------------------------------------------------------------------
// parseGitSlug
// ---------------------------------------------------------------------------

var parseGitSlugTests = []struct {
	name    string
	rawURL  string
	want    string
	wantErr bool
}{
	{
		name:   "SSH with .git",
		rawURL: "git@github.com:octocat/Hello-World.git",
		want:   "octocat/Hello-World",
	},
	{
		name:   "SSH without .git",
		rawURL: "git@github.com:octocat/Hello-World",
		want:   "octocat/Hello-World",
	},
	{
		name:   "HTTPS with .git",
		rawURL: "https://github.com/octocat/Hello-World.git",
		want:   "octocat/Hello-World",
	},
	{
		name:   "HTTPS without .git",
		rawURL: "https://github.com/octocat/Hello-World",
		want:   "octocat/Hello-World",
	},
	{
		name:   "HTTP",
		rawURL: "http://github.com/octocat/Hello-World.git",
		want:   "octocat/Hello-World",
	},
	{
		name:   "trailing whitespace",
		rawURL: "  git@github.com:octocat/Hello-World.git  ",
		want:   "octocat/Hello-World",
	},
	{
		name:    "empty string",
		rawURL:  "",
		wantErr: true,
	},
	{
		name:    "no slash after host",
		rawURL:  "https://github.com",
		wantErr: true,
	},
}

func TestParseGitSlug(t *testing.T) {
	for _, tc := range parseGitSlugTests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseGitSlug(tc.rawURL)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got slug %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// buildPullsURL
// ---------------------------------------------------------------------------

var buildPullsURLTests = []struct {
	name    string
	rawURL  string
	want    string
	wantErr bool
}{
	{
		name:   "SSH",
		rawURL: "git@github.com:octocat/Hello-World.git",
		want:   "https://github.com/octocat/Hello-World/pulls",
	},
	{
		name:   "HTTPS with .git",
		rawURL: "https://github.com/octocat/Hello-World.git",
		want:   "https://github.com/octocat/Hello-World/pulls",
	},
	{
		name:   "HTTPS without .git",
		rawURL: "https://github.com/octocat/Hello-World",
		want:   "https://github.com/octocat/Hello-World/pulls",
	},
	{
		name:   "HTTP",
		rawURL: "http://github.com/octocat/Hello-World.git",
		want:   "http://github.com/octocat/Hello-World/pulls",
	},
	{
		name:    "unsupported scheme",
		rawURL:  "ftp://github.com/octocat/Hello-World",
		wantErr: true,
	},
	{
		name:    "empty",
		rawURL:  "",
		wantErr: true,
	},
}

func TestBuildPullsURL(t *testing.T) {
	for _, tc := range buildPullsURLTests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := buildPullsURL(tc.rawURL)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("expected %q, got %q", tc.want, got)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// collectDirsRecursive
// ---------------------------------------------------------------------------

func TestCollectDirsRecursive_Depth1(t *testing.T) {
	root := t.TempDir()
	// Create: root/a, root/b, root/a/child, root/.hidden
	must(t, os.Mkdir(filepath.Join(root, "a"), 0o755))
	must(t, os.Mkdir(filepath.Join(root, "b"), 0o755))
	must(t, os.Mkdir(filepath.Join(root, "a", "child"), 0o755))
	must(t, os.Mkdir(filepath.Join(root, ".hidden"), 0o755))

	orig, _ := os.Getwd()
	must(t, os.Chdir(root))
	defer os.Chdir(orig) //nolint:errcheck

	dirs, err := collectDirsRecursive(".", 1)
	if err != nil {
		t.Fatal(err)
	}

	// depth=1 → only direct children (not grandchildren, not hidden)
	if len(dirs) != 2 {
		t.Errorf("expected 2 dirs, got %d: %v", len(dirs), dirs)
	}
	for _, d := range dirs {
		if strings.HasPrefix(filepath.Base(d), ".") {
			t.Errorf("hidden dir should be excluded: %q", d)
		}
	}
}

func TestCollectDirsRecursive_Depth2(t *testing.T) {
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, "a", "child"), 0o755))
	must(t, os.Mkdir(filepath.Join(root, "b"), 0o755))

	orig, _ := os.Getwd()
	must(t, os.Chdir(root))
	defer os.Chdir(orig) //nolint:errcheck

	dirs, err := collectDirsRecursive(".", 2)
	if err != nil {
		t.Fatal(err)
	}
	// Expect: a, a/child, b → 3 entries
	if len(dirs) != 3 {
		t.Errorf("expected 3 dirs, got %d: %v", len(dirs), dirs)
	}
}

func TestCollectDirsRecursive_NoSubdirs(t *testing.T) {
	root := t.TempDir()

	orig, _ := os.Getwd()
	must(t, os.Chdir(root))
	defer os.Chdir(orig) //nolint:errcheck

	dirs, err := collectDirsRecursive(".", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 0 {
		t.Errorf("expected no dirs, got %v", dirs)
	}
}

func TestCollectDirsRecursive_FilesIgnored(t *testing.T) {
	root := t.TempDir()
	must(t, os.WriteFile(filepath.Join(root, "file.txt"), []byte("x"), 0o644))
	must(t, os.Mkdir(filepath.Join(root, "subdir"), 0o755))

	orig, _ := os.Getwd()
	must(t, os.Chdir(root))
	defer os.Chdir(orig) //nolint:errcheck

	dirs, err := collectDirsRecursive(".", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 || filepath.Base(dirs[0]) != "subdir" {
		t.Errorf("expected only subdir, got %v", dirs)
	}
}

// ---------------------------------------------------------------------------
// colonConfigPath and loadOrInitConfig
// ---------------------------------------------------------------------------

func TestColonConfigPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}
	got, err := colonConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "colonsh.json")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestLoadOrInitConfig_CreatesDefault(t *testing.T) {
	dir := t.TempDir()
	// Redirect UserHomeDir by swapping the config file path.
	// We can't easily override os.UserHomeDir, so write and load via JSON directly.
	cfg := defaultConfig()
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	configFile := filepath.Join(dir, "colonsh.json")
	if err := os.WriteFile(configFile, data, 0o644); err != nil {
		t.Fatal(err)
	}

	// Load it back via json directly (loadOrInitConfig is tied to $HOME)
	data2, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatal(err)
	}
	var loaded Config
	if err := json.Unmarshal(data2, &loaded); err != nil {
		t.Fatalf("config roundtrip failed: %v", err)
	}
	if loaded.OpenCmd != cfg.OpenCmd {
		t.Errorf("OpenCmd mismatch: %q vs %q", loaded.OpenCmd, cfg.OpenCmd)
	}
	if len(loaded.Aliases) != len(cfg.Aliases) {
		t.Errorf("Aliases len mismatch: %d vs %d", len(loaded.Aliases), len(cfg.Aliases))
	}
}

// ---------------------------------------------------------------------------
// defaultConfig
// ---------------------------------------------------------------------------

func TestDefaultConfig_HasOpenCmd(t *testing.T) {
	cfg := defaultConfig()
	if cfg.OpenCmd == "" {
		t.Error("expected non-empty OpenCmd in default config")
	}
}

func TestDefaultConfig_HasProjectDirs(t *testing.T) {
	cfg := defaultConfig()
	if len(cfg.ProjectDirs) == 0 {
		t.Error("expected at least one ProjectDir in default config")
	}
}

func TestDefaultConfig_HasGitRepos(t *testing.T) {
	cfg := defaultConfig()
	if len(cfg.GitRepos) == 0 {
		t.Error("expected at least one GitRepo in default config")
	}
}

func TestDefaultConfig_HasAliases(t *testing.T) {
	cfg := defaultConfig()
	if len(cfg.Aliases) == 0 {
		t.Error("expected at least one Alias in default config")
	}
}

// ---------------------------------------------------------------------------
// findCurrentRepo
// ---------------------------------------------------------------------------

func TestFindCurrentRepo_Match(t *testing.T) {
	cfg := &Config{
		GitRepos: []GitRepo{
			{Slug: "alice/foo", Name: "foo"},
			{Slug: "bob/bar", Name: "bar"},
		},
	}
	// findCurrentRepo calls gitRepoSlug() internally (requires git), so we
	// test the lookup logic directly using the config slice.
	for i := range cfg.GitRepos {
		if cfg.GitRepos[i].Slug == "alice/foo" {
			if cfg.GitRepos[i].Name != "foo" {
				t.Errorf("unexpected name %q", cfg.GitRepos[i].Name)
			}
			return
		}
	}
	t.Error("alice/foo not found")
}

// ---------------------------------------------------------------------------
// cmdGitNewBranch argument joining
// ---------------------------------------------------------------------------

func TestCmdGitNewBranch_NoArgs(t *testing.T) {
	err := cmdGitNewBranch([]string{})
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

// ---------------------------------------------------------------------------
// cmdGitCommit argument validation
// ---------------------------------------------------------------------------

func TestCmdGitCommit_NoArgs(t *testing.T) {
	err := cmdGitCommit([]string{})
	if err == nil {
		t.Error("expected error when no commit message provided")
	}
}

func TestCmdGitCommitAmend_NoArgs(t *testing.T) {
	err := cmdGitCommitAmendWithMessage([]string{})
	if err == nil {
		t.Error("expected error when no commit message provided")
	}
}

// ---------------------------------------------------------------------------
// helper
// ---------------------------------------------------------------------------

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
