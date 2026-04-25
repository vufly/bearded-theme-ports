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
