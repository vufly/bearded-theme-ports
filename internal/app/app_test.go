package app

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestParseBuildOptionsDefaultInstallCoversInstallerScripts(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "..", "scripts"))
	if err != nil {
		t.Fatal(err)
	}
	scriptTargets := make(map[string]bool)
	for _, entry := range entries {
		name := entry.Name()
		ext := filepath.Ext(name)
		if strings.HasPrefix(name, "install-") && (ext == ".sh" || ext == ".ps1") {
			scriptTargets[strings.TrimSuffix(strings.TrimPrefix(name, "install-"), ext)] = true
		}
	}
	if !scriptTargets["bat"] || !scriptTargets["codex"] {
		t.Fatal("expected install-only consumer scripts for bat and codex")
	}

	for _, flag := range []string{"--install", "-i"} {
		t.Run(flag, func(t *testing.T) {
			options, err := parseBuildOptions([]string{flag})
			if err != nil {
				t.Fatal(err)
			}
			if !options.install || !slices.Equal(options.targets, allTargets()) {
				t.Fatalf("default build options = %#v, want all build targets with install enabled", options)
			}
			if len(options.installTargets) != len(scriptTargets) || !slices.IsSorted(options.installTargets) {
				t.Fatalf("install targets = %v, want %d sorted installer targets", options.installTargets, len(scriptTargets))
			}
			for target := range scriptTargets {
				if !slices.Contains(options.installTargets, target) {
					t.Errorf("default install missing %q", target)
				}
			}
		})
	}
}

func TestParseBuildOptionsDefaultBuildDoesNotInstall(t *testing.T) {
	options, err := parseBuildOptions(nil)
	if err != nil {
		t.Fatal(err)
	}
	if options.install || len(options.installTargets) != 0 || !slices.Equal(options.targets, allTargets()) {
		t.Fatalf("default build options = %#v, want all build targets without installers", options)
	}
}

func TestParseBuildOptionsExplicitInstallDeduplicatesSharedBuild(t *testing.T) {
	options, err := parseBuildOptions([]string{"--install", "bat", "codex", "tmtheme", "bat", "wezterm"})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(options.targets, []string{"tmtheme", "wezterm"}) {
		t.Fatalf("build targets = %v, want tmtheme and wezterm once each", options.targets)
	}
	if !slices.Equal(options.installTargets, []string{"bat", "codex", "wezterm"}) {
		t.Fatalf("install targets = %v, want only selected consumers and wezterm", options.installTargets)
	}
}

func TestParseTargetsInstallBatUsesTMThemeBuild(t *testing.T) {
	targets, installTargets, err := parseTargets([]string{"bat"}, true)
	if err != nil {
		t.Fatalf("parseTargets() error = %v", err)
	}

	if len(targets) != 1 || targets[0] != "tmtheme" {
		t.Fatalf("parseTargets() targets = %#v, want []string{\"tmtheme\"}", targets)
	}
	if len(installTargets) != 1 || installTargets[0] != "bat" {
		t.Fatalf("parseTargets() installTargets = %#v, want []string{\"bat\"}", installTargets)
	}
}

func TestParseTargetsRejectsBatWithoutInstall(t *testing.T) {
	_, _, err := parseTargets([]string{"bat"}, false)
	if err == nil {
		t.Fatal("parseTargets() error = nil, want error")
	}
}

func TestParseTargetsInstallCodexUsesTMThemeBuild(t *testing.T) {
	targets, installTargets, err := parseTargets([]string{"codex"}, true)
	if err != nil {
		t.Fatalf("parseTargets() error = %v", err)
	}

	if len(targets) != 1 || targets[0] != "tmtheme" {
		t.Fatalf("parseTargets() targets = %#v, want []string{\"tmtheme\"}", targets)
	}
	if len(installTargets) != 1 || installTargets[0] != "codex" {
		t.Fatalf("parseTargets() installTargets = %#v, want []string{\"codex\"}", installTargets)
	}
}

func TestParseTargetsRejectsCodexWithoutInstall(t *testing.T) {
	_, _, err := parseTargets([]string{"codex"}, false)
	if err == nil {
		t.Fatal("parseTargets() error = nil, want error")
	}
}
