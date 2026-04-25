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
	case "bat", "codex", "helix", "neovim", "opencode", "wezterm":
		return true
	default:
		return false
	}
}

func Install(root string, target string) (string, error) {
	var sourceDir string
	var targetDir string
	var postInstall func() error

	switch target {
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
	case "helix":
		sourceDir = source.HelixOutputDir(root)
		targetDir = helixThemesDir()
	case "neovim":
		sourceDir = source.NeovimOutputDir(root)
		targetDir = neovimColorsDir()
	case "opencode":
		sourceDir = source.OpenCodeOutputDir(root)
		targetDir = opencodeThemesDir()
	case "wezterm":
		sourceDir = source.WezTermOutputDir(root)
		targetDir = weztermThemesDir()
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
