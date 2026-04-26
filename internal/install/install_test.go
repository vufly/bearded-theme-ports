package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallBatCopiesThemesAndBuildsCache(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "dist", "tmtheme")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "bearded-theme.tmTheme"), []byte("plist"), 0o644); err != nil {
		t.Fatalf("WriteFile() theme error = %v", err)
	}

	batConfigDir := filepath.Join(root, "bat-config")
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() bin error = %v", err)
	}

	logPath := filepath.Join(root, "bat.log")
	batScript := "#!/bin/sh\n" +
		"echo \"$@\" >> " + shellQuote(logPath) + "\n" +
		"if [ \"$1\" = \"--config-dir\" ]; then\n" +
		"  printf %s " + shellQuote(batConfigDir) + "\n" +
		"  exit 0\n" +
		"fi\n" +
		"exit 0\n"
	batPath := filepath.Join(binDir, "bat")
	if err := os.WriteFile(batPath, []byte(batScript), 0o755); err != nil {
		t.Fatalf("WriteFile() bat script error = %v", err)
	}

	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	targetDir, err := Install(root, "bat")
	if err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if got, want := targetDir, filepath.Join(batConfigDir, "themes"); got != want {
		t.Fatalf("Install() targetDir = %q, want %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(targetDir, "bearded-theme.tmTheme")); err != nil {
		t.Fatalf("installed theme missing: %v", err)
	}

	logContent, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile() log error = %v", err)
	}
	output := string(logContent)
	if !strings.Contains(output, "--config-dir") || !strings.Contains(output, "cache --build") {
		t.Fatalf("expected bat invocations in log, got:\n%s", output)
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
}

func TestSupportedTarget_CoversInstallScriptTargets(t *testing.T) {
	supported := []string{
		"alacritty", "bat", "codex", "delta", "ghostty",
		"helix", "kitty", "neovim", "opencode", "termux",
		"wezterm", "zellij",
	}
	for _, target := range supported {
		if !SupportedTarget(target) {
			t.Errorf("SupportedTarget(%q) = false, want true", target)
		}
	}

	// These are intentionally not installable; if anyone re-adds them they
	// must do so consciously and update the comment in install.go.
	notInstallable := []string{"firefox-color", "windows-terminal", "lazygit", "tmtheme"}
	for _, target := range notInstallable {
		if SupportedTarget(target) {
			t.Errorf("SupportedTarget(%q) = true, want false", target)
		}
	}
}

func TestInstall_SimpleTargets_CopyDirContentsIntoConfigSubdir(t *testing.T) {
	cases := []struct {
		target     string
		distSub    string
		fileName   string
		configSub  string
	}{
		{target: "alacritty", distSub: "alacritty", fileName: "bearded-theme-monokai-stone.toml", configSub: "alacritty/themes"},
		{target: "ghostty", distSub: "ghostty", fileName: "bearded-theme-monokai-stone", configSub: "ghostty/themes"},
		{target: "kitty", distSub: "kitty", fileName: "bearded-theme-monokai-stone.conf", configSub: "kitty/themes"},
		{target: "zellij", distSub: "zellij", fileName: "bearded-theme-monokai-stone.kdl", configSub: "zellij/themes"},
	}

	for _, tc := range cases {
		t.Run(tc.target, func(t *testing.T) {
			root := t.TempDir()
			home := t.TempDir()

			// Drop the source artifact the installer expects to find.
			distDir := filepath.Join(root, "dist", tc.distSub)
			if err := os.MkdirAll(distDir, 0o755); err != nil {
				t.Fatalf("MkdirAll() dist error = %v", err)
			}
			if err := os.WriteFile(filepath.Join(distDir, tc.fileName), []byte("payload"), 0o644); err != nil {
				t.Fatalf("WriteFile() dist error = %v", err)
			}

			// Force the installer to land under our temp HOME by clearing
			// XDG_CONFIG_HOME and pointing HOME at the tempdir.
			t.Setenv("HOME", home)
			t.Setenv("XDG_CONFIG_HOME", "")

			targetDir, err := Install(root, tc.target)
			if err != nil {
				t.Fatalf("Install(%q) error = %v", tc.target, err)
			}

			wantDir := filepath.Join(home, ".config", tc.configSub)
			if targetDir != wantDir {
				t.Fatalf("Install(%q) targetDir = %q, want %q", tc.target, targetDir, wantDir)
			}
			if _, err := os.Stat(filepath.Join(targetDir, tc.fileName)); err != nil {
				t.Fatalf("installed file missing for %q: %v", tc.target, err)
			}
		})
	}
}

func TestInstall_Termux_ReplacesColorsPropertiesAndBacksUpExisting(t *testing.T) {
	root := t.TempDir()
	home := t.TempDir()

	// Two upstream variants so we can confirm slug-aware selection.
	distDir := filepath.Join(root, "dist", "termux")
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	defaultPath := filepath.Join(distDir, "bearded-theme-monokai-stone.properties")
	customPath := filepath.Join(distDir, "bearded-theme-vivid-purple.properties")
	if err := os.WriteFile(defaultPath, []byte("default-payload"), 0o644); err != nil {
		t.Fatalf("WriteFile() default error = %v", err)
	}
	if err := os.WriteFile(customPath, []byte("custom-payload"), 0o644); err != nil {
		t.Fatalf("WriteFile() custom error = %v", err)
	}

	// Pre-existing colors.properties so we exercise the backup path.
	termuxDir := filepath.Join(home, ".termux")
	if err := os.MkdirAll(termuxDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() termux error = %v", err)
	}
	existingPath := filepath.Join(termuxDir, "colors.properties")
	if err := os.WriteFile(existingPath, []byte("existing-payload"), 0o644); err != nil {
		t.Fatalf("WriteFile() existing error = %v", err)
	}

	t.Setenv("HOME", home)
	// Empty PATH ensures termux-reload-settings is treated as absent —
	// keeps the test hermetic on every host.
	t.Setenv("PATH", "")

	// First install: default slug, should back up the existing file.
	targetPath, err := Install(root, "termux")
	if err != nil {
		t.Fatalf("Install(termux) default error = %v", err)
	}
	if targetPath != existingPath {
		t.Fatalf("Install(termux) targetPath = %q, want %q", targetPath, existingPath)
	}
	if got, _ := os.ReadFile(existingPath); string(got) != "default-payload" {
		t.Fatalf("colors.properties contents = %q, want default-payload", got)
	}
	if got, _ := os.ReadFile(existingPath + ".bak"); string(got) != "existing-payload" {
		t.Fatalf("backup contents = %q, want existing-payload", got)
	}

	// Second install: custom slug via TERMUX_THEME, backup must NOT be
	// overwritten (it captures the user's pre-Bearded scheme).
	t.Setenv("TERMUX_THEME", "bearded-theme-vivid-purple")
	if _, err := Install(root, "termux"); err != nil {
		t.Fatalf("Install(termux) custom error = %v", err)
	}
	if got, _ := os.ReadFile(existingPath); string(got) != "custom-payload" {
		t.Fatalf("colors.properties contents = %q, want custom-payload", got)
	}
	if got, _ := os.ReadFile(existingPath + ".bak"); string(got) != "existing-payload" {
		t.Fatalf("backup overwritten on second install: got %q", got)
	}
}

func TestInstall_Termux_UnknownSlugErrors(t *testing.T) {
	root := t.TempDir()
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "dist", "termux"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	t.Setenv("HOME", home)
	t.Setenv("TERMUX_THEME", "bearded-theme-does-not-exist")

	if _, err := Install(root, "termux"); err == nil {
		t.Fatal("Install(termux) with unknown slug should error")
	}
}

func TestInstall_Delta_CopiesAllGitconfigFiles(t *testing.T) {
	root := t.TempDir()
	xdg := t.TempDir()

	// Source assets: consolidated file plus one single-theme file.
	deltaDir := filepath.Join(root, "dist", "delta")
	if err := os.MkdirAll(deltaDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(deltaDir, "bearded-theme.gitconfig"), []byte("[delta \"x\"]\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(deltaDir, "bearded-theme-monokai-stone.gitconfig"), []byte("[delta \"stone\"]\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() theme error = %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", xdg)

	targetDir, err := Install(root, "delta")
	if err != nil {
		t.Fatalf("Install(delta) error = %v", err)
	}

	wantDir := filepath.Join(xdg, "git")
	if targetDir != wantDir {
		t.Fatalf("Install(delta) targetDir = %q, want %q", targetDir, wantDir)
	}
	for _, name := range []string{"bearded-theme.gitconfig", "bearded-theme-monokai-stone.gitconfig"} {
		path := filepath.Join(wantDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("gitconfig missing at target: %v", err)
		}
	}
	if got, err := os.ReadFile(filepath.Join(wantDir, "bearded-theme.gitconfig")); err != nil {
		t.Fatalf("ReadFile() installed gitconfig error = %v", err)
	} else if string(got) != "[delta \"x\"]\n" {
		t.Fatalf("installed gitconfig = %q, want source contents", got)
	}
	if got, err := os.ReadFile(filepath.Join(wantDir, "bearded-theme-monokai-stone.gitconfig")); err != nil {
		t.Fatalf("ReadFile() installed per-theme gitconfig error = %v", err)
	} else if string(got) != "[delta \"stone\"]\n" {
		t.Fatalf("installed per-theme gitconfig = %q, want source contents", got)
	}
}
