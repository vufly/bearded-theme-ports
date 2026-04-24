package app

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/output"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/targets/tmtheme"
	"bearded-theme-ports/internal/targets/wezterm"
)

type builderFunc func(root string, themes []model.ThemeFile) ([]string, error)

func Run(args []string) error {
	root, err := source.FindRepoRoot()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		printUsage()
		return nil
	case "sync":
		return source.Sync(root)
	case "prepare-upstream":
		return source.PrepareUpstream(root)
	case "prepare-and-build", "build-all":
		return prepareAndBuild(root, args[1:])
	case "build":
		return build(root, args[1:])
	case "list":
		return list(root, args[1:])
	case "doctor":
		return doctor(root)
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], usageText)
	}
}

func build(root string, args []string) error {
	targets, err := parseTargets(args)
	if err != nil {
		return err
	}

	themes, err := source.LoadThemes(root)
	if err != nil {
		return err
	}

	for _, target := range targets {
		if _, err := builders[target](root, themes); err != nil {
			return err
		}
	}

	commitSHA, err := source.UpstreamCommitSHA(root)
	if err != nil {
		return err
	}

	if err := output.WriteMetadata(root, output.Metadata{
		GeneratedTargets: targets,
		InputThemeCount:  len(themes),
		SourcePath:       source.InputThemesDir(root),
		UpstreamCommit:   commitSHA,
		UpstreamRepoURL:  source.UpstreamRepoURL,
	}); err != nil {
		return err
	}

	fmt.Printf("✅ Generated %d themes for targets: %s\n", len(themes), formatTargets(targets))
	return nil
}

func formatTargets(targets []string) string {
	formatted := make([]string, 0, len(targets))
	for _, target := range targets {
		formatted = append(formatted, fmt.Sprintf("\033[1;33m%s\033[0m", target))
	}

	return strings.Join(formatted, ", ")
}

func prepareAndBuild(root string, args []string) error {
	if err := source.Sync(root); err != nil {
		return err
	}

	if err := source.PrepareUpstream(root); err != nil {
		return err
	}

	return build(root, args)
}

func parseTargets(args []string) ([]string, error) {
	if len(args) == 0 {
		return allTargets(), nil
	}

	seen := make(map[string]bool, len(args))
	targets := make([]string, 0, len(args))
	for _, target := range args {
		if _, ok := builders[target]; !ok {
			return nil, fmt.Errorf("unsupported build target %q", target)
		}

		if seen[target] {
			continue
		}

		seen[target] = true
		targets = append(targets, target)
	}

	return targets, nil
}

func allTargets() []string {
	targets := make([]string, 0, len(builders))
	for target := range builders {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	return targets
}

var builders = map[string]builderFunc{
	"tmtheme": tmtheme.Build,
	"wezterm": wezterm.Build,
}

func list(root string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("unsupported list command\n\n%s", usageText)
	}

	switch args[0] {
	case "themes":
		return listThemes(root)
	case "targets":
		for _, target := range allTargets() {
			fmt.Println(target)
		}
		return nil
	default:
		return fmt.Errorf("unsupported list command\n\n%s", usageText)
	}
}

func listThemes(root string) error {
	themes, err := source.LoadThemes(root)
	if err != nil {
		return err
	}

	for _, theme := range themes {
		fmt.Printf("%s\t%s\n", theme.Slug, filepathBase(theme.Path))
	}

	return nil
}

func doctor(root string) error {
	issues := make([]string, 0, 4)

	if err := source.CheckExecutable("git"); err != nil {
		issues = append(issues, err.Error())
	}

	if _, err := source.DetectPackageManager(root); err != nil {
		issues = append(issues, err.Error())
	}

	if _, err := os.Stat(source.UpstreamDir(root)); err != nil {
		issues = append(issues, fmt.Sprintf("upstream repo missing at %s", source.UpstreamDir(root)))
	}

	if _, err := os.Stat(source.InputThemesDir(root)); err != nil {
		issues = append(issues, fmt.Sprintf("prepared theme JSON missing at %s", source.InputThemesDir(root)))
	}

	if len(issues) > 0 {
		return errors.New(strings.Join(issues, "\n"))
	}

	fmt.Println("doctor: ok")
	return nil
}

func filepathBase(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	return parts[len(parts)-1]
}

const usageText = `Usage:
  bearded-theme-ports sync
  bearded-theme-ports prepare-upstream
  bearded-theme-ports prepare-and-build [targets...]
  bearded-theme-ports build-all [targets...]
  bearded-theme-ports build [targets...]
  bearded-theme-ports list themes
  bearded-theme-ports list targets
  bearded-theme-ports doctor`

func printUsage() {
	fmt.Println(usageText)
}
