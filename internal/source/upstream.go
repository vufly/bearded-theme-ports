package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"bearded-theme-ports/internal/model"
)

const UpstreamRepoURL = "https://github.com/BeardedBear/bearded-theme"

func FindRepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	current := wd
	for {
		if fileExists(filepath.Join(current, ".git")) {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", errors.New("could not locate repository root")
		}

		current = parent
	}
}

func UpstreamDir(root string) string {
	return filepath.Join(root, ".cache", "upstream", "bearded-theme")
}

func InputThemesDir(root string) string {
	return filepath.Join(UpstreamDir(root), "dist", "vscode", "themes")
}

func WezTermOutputDir(root string) string {
	return filepath.Join(root, "dist", "wezterm")
}

func LegacyTargetTypesDir(root string) string {
	return filepath.Join(root, "dist", "terminals")
}

func MetadataDir(root string) string {
	return filepath.Join(root, "dist", "metadata")
}

func Sync(root string) error {
	cacheParent := filepath.Dir(UpstreamDir(root))
	if err := os.MkdirAll(cacheParent, 0o755); err != nil {
		return err
	}

	if fileExists(filepath.Join(UpstreamDir(root), ".git")) {
		return runCommand(root, "git", "-C", UpstreamDir(root), "pull", "--ff-only")
	}

	return runCommand(root, "git", "clone", "--depth=1", UpstreamRepoURL, UpstreamDir(root))
}

func PrepareUpstream(root string) error {
	if !fileExists(filepath.Join(UpstreamDir(root), ".git")) {
		return fmt.Errorf("upstream repo missing at %s; run sync first", UpstreamDir(root))
	}

	if err := runCommand(UpstreamDir(root), "npm", "ci"); err != nil {
		return err
	}

	return runCommand(UpstreamDir(root), "npm", "run", "build:vscode")
}

func UpstreamCommitSHA(root string) (string, error) {
	if !fileExists(filepath.Join(UpstreamDir(root), ".git")) {
		return "", fmt.Errorf("upstream repo missing at %s", UpstreamDir(root))
	}

	output, err := exec.Command("git", "-C", UpstreamDir(root), "rev-parse", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func LoadThemes(root string) ([]model.ThemeFile, error) {
	patterns, err := filepath.Glob(filepath.Join(InputThemesDir(root), "bearded-theme-*.json"))
	if err != nil {
		return nil, err
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf("no theme JSON files found in %s; run prepare-upstream first", InputThemesDir(root))
	}

	sort.Strings(patterns)
	themes := make([]model.ThemeFile, 0, len(patterns))

	for _, path := range patterns {
		fileName := filepath.Base(path)
		if fileName == "ui-key-tester.json" {
			continue
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var theme model.VSCodeTheme
		if err := json.Unmarshal(content, &theme); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}

		themes = append(themes, model.ThemeFile{
			Path:  path,
			Slug:  strings.TrimSuffix(fileName, ".json"),
			Theme: theme,
		})
	}

	return themes, nil
}

func CheckExecutable(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("missing executable %q", name)
	}

	return nil
}

func runCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
