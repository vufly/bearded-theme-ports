package install

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"bearded-theme-ports/internal/source"
)

func SupportedTarget(target string) bool {
	switch target {
	case "alacritty",
		"bat",
		"codex",
		"delta",
		"ghostty",
		"helix",
		"kitty",
		"neovim",
		"opencode",
		"termux",
		"wezterm",
		"zellij":
		return true
	}
	// Intentionally not installable:
	//   - firefox-color: URL-based (no filesystem destination)
	//   - windows-terminal: requires merging schemes into settings.json
	//   - lazygit: requires merging YAML under gui.theme in config.yml
	//   - tmtheme: consumed by the bat install path
	return false
}

func Install(root string, target string) (string, error) {
	// termux and delta don't follow the "copy a directory of files into
	// ~/.config/<tool>/themes/" pattern, so they get dedicated installers.
	switch target {
	case "termux":
		return installTermux(root)
	case "delta":
		return installDelta(root)
	}

	var sourceDir string
	var targetDir string
	var postInstall func() error

	switch target {
	case "alacritty":
		sourceDir = source.AlacrittyOutputDir(root)
		targetDir = alacrittyThemesDir()
	case "bat":
		batBin, batTargetDir, err := batThemesDir()
		if err != nil {
			return "", err
		}
		sourceDir = source.TMThemeOutputDir(root)
		targetDir = batTargetDir
		postInstall = func() error {
			return runBatCacheBuild(batBin)
		}
	case "codex":
		sourceDir = source.CodexOutputDir(root)
		targetDir = codexThemesDir()
	case "ghostty":
		sourceDir = source.GhosttyOutputDir(root)
		targetDir = ghosttyThemesDir()
	case "helix":
		sourceDir = source.HelixOutputDir(root)
		targetDir = helixThemesDir()
	case "kitty":
		sourceDir = source.KittyOutputDir(root)
		targetDir = kittyThemesDir()
	case "neovim":
		sourceDir = source.NeovimOutputDir(root)
		targetDir = neovimColorsDir()
	case "opencode":
		sourceDir = source.OpenCodeOutputDir(root)
		targetDir = opencodeThemesDir()
	case "wezterm":
		sourceDir = source.WezTermOutputDir(root)
		targetDir = weztermThemesDir()
	case "zellij":
		sourceDir = source.ZellijOutputDir(root)
		targetDir = zellijThemesDir()
	default:
		return "", fmt.Errorf("install not supported for target %q", target)
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	if err := copyDirContents(sourceDir, targetDir); err != nil {
		return "", err
	}
	if postInstall != nil {
		if err := postInstall(); err != nil {
			return "", err
		}
	}

	return targetDir, nil
}

func copyDirContents(sourceDir string, targetDir string) error {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(sourceDir, entry.Name())
		targetPath := filepath.Join(targetDir, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(targetPath, 0o755); err != nil {
				return err
			}
			if err := copyDirContents(sourcePath, targetPath); err != nil {
				return err
			}
			continue
		}

		if err := copyFile(sourcePath, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(sourcePath string, targetPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return err
	}

	return targetFile.Close()
}

func configRootDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return filepath.Join(homeDir, ".config")
}

func helixThemesDir() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "helix", "themes")
		}
	}

	return filepath.Join(configRootDir(), "helix", "themes")
}

func neovimColorsDir() string {
	if runtime.GOOS == "windows" {
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "nvim", "colors")
		}
	}

	return filepath.Join(configRootDir(), "nvim", "colors")
}

func weztermThemesDir() string {
	return filepath.Join(configRootDir(), "wezterm", "themes", "bearded-theme")
}

func codexThemesDir() string {
	if codexHome := os.Getenv("CODEX_HOME"); codexHome != "" {
		return filepath.Join(codexHome, "themes")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".codex", "themes")
	}

	return filepath.Join(homeDir, ".codex", "themes")
}

func opencodeThemesDir() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "opencode", "themes")
		}
	}

	return filepath.Join(configRootDir(), "opencode", "themes")
}

func kittyThemesDir() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "kitty", "themes")
		}
	}

	return filepath.Join(configRootDir(), "kitty", "themes")
}

func alacrittyThemesDir() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "alacritty", "themes")
		}
	}

	return filepath.Join(configRootDir(), "alacritty", "themes")
}

func ghosttyThemesDir() string {
	// Ghostty has no Windows port; honor XDG/HOME on every platform.
	return filepath.Join(configRootDir(), "ghostty", "themes")
}

func zellijThemesDir() string {
	// Use $HOME/.config/zellij/themes on every platform (honoring
	// XDG_CONFIG_HOME when set) so the install path stays consistent
	// across darwin, linux and Windows.
	return filepath.Join(configRootDir(), "zellij", "themes")
}

// installTermux copies a single chosen .properties file over
// ~/.termux/colors.properties (Termux only renders one active scheme at a
// time). The slug defaults to bearded-theme-monokai-stone but can be
// overridden via the TERMUX_THEME env var so users can pick any variant
// without changing the install command shape.
func installTermux(root string) (string, error) {
	slug := strings.TrimSpace(os.Getenv("TERMUX_THEME"))
	if slug == "" {
		slug = "bearded-theme-monokai-stone"
	}

	sourcePath := filepath.Join(source.TermuxOutputDir(root), slug+".properties")
	if _, err := os.Stat(sourcePath); err != nil {
		return "", fmt.Errorf("termux theme %q not found at %s: %w", slug, sourcePath, err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	targetDir := filepath.Join(homeDir, ".termux")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}
	targetPath := filepath.Join(targetDir, "colors.properties")

	// Save a one-shot backup of any existing scheme so users can roll back
	// after a first overwrite.
	if _, err := os.Stat(targetPath); err == nil {
		backupPath := targetPath + ".bak"
		if _, backupErr := os.Stat(backupPath); os.IsNotExist(backupErr) {
			if err := copyFile(targetPath, backupPath); err != nil {
				return "", err
			}
		}
	}

	if err := copyFile(sourcePath, targetPath); err != nil {
		return "", err
	}

	// Refresh Termux if we're actually inside it; ignore errors so this
	// works fine when invoked from a regular Linux/macOS host (e.g. CI).
	if termuxBin, err := exec.LookPath("termux-reload-settings"); err == nil {
		_ = exec.Command(termuxBin).Run()
	}

	return targetPath, nil
}

// installDelta drops the consolidated bearded-theme.gitconfig into the
// user's git config dir and registers it via `git config --global --add
// include.path`, matching what scripts/install-delta.sh does. Re-running
// is idempotent: an existing matching include.path entry is preserved
// rather than duplicated.
func installDelta(root string) (string, error) {
	sourcePath := filepath.Join(source.DeltaOutputDir(root), "bearded-theme.gitconfig")
	if _, err := os.Stat(sourcePath); err != nil {
		return "", fmt.Errorf("delta consolidated gitconfig not found at %s: %w", sourcePath, err)
	}

	targetDir := filepath.Join(configRootDir(), "git")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}
	targetPath := filepath.Join(targetDir, "bearded-theme.gitconfig")
	if err := copyFile(sourcePath, targetPath); err != nil {
		return "", err
	}

	if err := registerDeltaInclude(targetPath); err != nil {
		return "", err
	}

	return targetPath, nil
}

func registerDeltaInclude(path string) error {
	gitBin, err := exec.LookPath("git")
	if err != nil {
		// git not installed: nothing to register, but the file is still in
		// place. Fall back to a manual instruction rather than failing.
		fmt.Fprintf(os.Stderr,
			"git not found in PATH; add this manually:\n  [include]\n      path = %s\n",
			path,
		)
		return nil
	}

	// `git config --get-all` exits non-zero when the key is unset; that's
	// not an error for our purposes — we just treat it as "nothing
	// registered yet".
	output, _ := exec.Command(gitBin, "config", "--global", "--get-all", "include.path").Output()
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if line == path {
			return nil
		}
	}

	cmd := exec.Command(gitBin, "config", "--global", "--add", "include.path", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func batThemesDir() (string, string, error) {
	batBin, err := resolveBatCommand()
	if err != nil {
		return "", "", err
	}
	if configDir := strings.TrimSpace(os.Getenv("BAT_CONFIG_DIR")); configDir != "" {
		return batBin, filepath.Join(configDir, "themes"), nil
	}

	output, err := exec.Command(batBin, "--config-dir").Output()
	if err != nil {
		return "", "", fmt.Errorf("resolve bat config dir: %w", err)
	}

	configDir := strings.TrimSpace(string(output))
	if configDir == "" {
		return "", "", fmt.Errorf("%s --config-dir returned an empty path", batBin)
	}

	return batBin, filepath.Join(configDir, "themes"), nil
}

func resolveBatCommand() (string, error) {
	for _, name := range []string{"bat", "batcat"} {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("missing bat executable: need bat or batcat")
}

func runBatCacheBuild(batBin string) error {
	cmd := exec.Command(batBin, "cache", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
