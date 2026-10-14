package install

import (
	"archive/zip"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const pigmentThemeFile = "bearded-theme-monokai-stone.tmTheme"

func TestInstallPiPigmentCopiesTMThemes(t *testing.T) {
	for _, custom := range []bool{false, true} {
		name := "default"
		if custom {
			name = "custom-agent-dir"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			targetDir := pigmentTestTargetDir(t, custom)
			sourceDir := filepath.Join(root, "dist", "tmtheme")
			if err := os.MkdirAll(sourceDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(sourceDir, pigmentThemeFile), []byte("plist"), 0o644); err != nil {
				t.Fatal(err)
			}
			gotDir, err := Install(root, "pi-pigment")
			if err != nil {
				t.Fatal(err)
			}
			if gotDir != targetDir {
				t.Fatalf("Install() directory = %q, want %q", gotDir, targetDir)
			}
			assertPigmentInstall(t, targetDir)
		})
	}
}

func TestInstallPiPigmentMissingTMThemes(t *testing.T) {
	pigmentTestTargetDir(t, false)
	if _, err := Install(t.TempDir(), "pi-pigment"); err == nil {
		t.Fatal("Install() should fail without tmtheme output")
	}
}

func TestPiPigmentInstallScripts(t *testing.T) {
	for _, ext := range []string{"sh", "ps1"} {
		for _, custom := range []bool{false, true} {
			name := ext + "/default"
			if custom {
				name = ext + "/custom-agent-dir"
			}
			t.Run(name, func(t *testing.T) {
				var shell string
				var err error
				if ext == "sh" {
					if runtime.GOOS == "windows" {
						t.Skip("POSIX installer")
					}
					shell, err = exec.LookPath("sh")
					if _, unzipErr := exec.LookPath("unzip"); unzipErr != nil {
						t.Skip("unzip unavailable")
					}
				} else {
					shell, err = exec.LookPath("pwsh")
					if err != nil {
						shell, err = exec.LookPath("powershell")
					}
				}
				if err != nil {
					t.Skip("installer shell unavailable")
				}

				targetDir := pigmentTestTargetDir(t, custom)
				root := t.TempDir()
				archivePath := filepath.Join(root, "themes.zip")
				writePigmentTestArchive(t, archivePath)
				t.Setenv("PIGMENT_TEST_ARCHIVE", archivePath)
				t.Setenv("PIGMENT_TEST_URL", "https://github.com/vufly/bearded-theme-ports/releases/latest/download/bearded-theme-ports-tmtheme.zip")
				scriptPath, err := filepath.Abs(filepath.Join("..", "..", "scripts", "install-pi-pigment."+ext))
				if err != nil {
					t.Fatal(err)
				}
				var cmd *exec.Cmd
				if ext == "sh" {
					binDir := filepath.Join(root, "bin")
					if err := os.MkdirAll(binDir, 0o755); err != nil {
						t.Fatal(err)
					}
					// Replace network access, but exercise real unzip and copy commands.
					mock := "#!/bin/sh\nset -eu\n[ \"$2\" = \"$PIGMENT_TEST_URL\" ]\ncp \"$PIGMENT_TEST_ARCHIVE\" \"$4\"\n"
					if err := os.WriteFile(filepath.Join(binDir, "curl"), []byte(mock), 0o755); err != nil {
						t.Fatal(err)
					}
					t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
					cmd = exec.Command(shell, scriptPath)
				} else {
					// PowerShell's HOME variable is initialized independently of HOME env.
					command := "Set-Variable -Name HOME -Value $env:HOME -Force; " +
						"function Invoke-WebRequest { param($Uri, $OutFile) " +
						"if ($Uri -ne $env:PIGMENT_TEST_URL) { throw 'Wrong release asset' }; " +
						"Copy-Item -LiteralPath $env:PIGMENT_TEST_ARCHIVE -Destination $OutFile }; " +
						"& '" + strings.ReplaceAll(scriptPath, "'", "''") + "'"
					cmd = exec.Command(shell, "-NoProfile", "-NonInteractive", "-Command", command)
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("installer failed: %v\n%s", err, output)
				}
				if !strings.Contains(string(output), "/pigment convert bearded-theme-monokai-stone") {
					t.Fatalf("missing conversion instructions: %s", output)
				}
				assertPigmentInstall(t, targetDir)
			})
		}
	}
}

func pigmentTestTargetDir(t *testing.T, custom bool) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("PI_CODING_AGENT_DIR", "")
	// Pigment uses Pi's home-based agent directory, not Windows/XDG config roots.
	t.Setenv("APPDATA", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	agentDir := filepath.Join(home, ".pi", "agent")
	if custom {
		agentDir = filepath.Join(home, "custom agent")
		t.Setenv("PI_CODING_AGENT_DIR", agentDir)
	}
	targetDir := filepath.Join(agentDir, "extensions", "pigment", "themes")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{
		filepath.Join(targetDir, "pigment-bearded-theme-monokai-stone.json"),
		filepath.Join(targetDir, "..", "config.jsonc"),
	} {
		if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return targetDir
}

func assertPigmentInstall(t *testing.T, targetDir string) {
	t.Helper()
	for path, want := range map[string]string{
		filepath.Join(targetDir, pigmentThemeFile):                           "plist",
		filepath.Join(targetDir, "pigment-bearded-theme-monokai-stone.json"): "existing",
		filepath.Join(targetDir, "..", "config.jsonc"):                       "existing",
	} {
		got, err := os.ReadFile(path)
		if err != nil || string(got) != want {
			t.Errorf("installed file %s = %q, %v; want %q", path, got, err, want)
		}
	}
}

func writePigmentTestArchive(t *testing.T, path string) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	archive := zip.NewWriter(file)
	entry, err := archive.Create(pigmentThemeFile)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := entry.Write([]byte("plist")); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
}
