package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"bearded-theme-ports/internal/source"
)

func SupportedTarget(target string) bool {
	switch target {
	case "helix", "neovim", "wezterm":
		return true
	default:
		return false
	}
}

func Install(root string, target string) (string, error) {
	var sourceDir string
	var targetDir string

	switch target {
	case "helix":
		sourceDir = source.HelixOutputDir(root)
		targetDir = helixThemesDir()
	case "neovim":
		sourceDir = source.NeovimOutputDir(root)
		targetDir = neovimColorsDir()
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
