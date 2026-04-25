package app

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"bearded-theme-ports/internal/install"
	"bearded-theme-ports/internal/model"
	"bearded-theme-ports/internal/output"
	"bearded-theme-ports/internal/source"
	"bearded-theme-ports/internal/targets/helix"
	"bearded-theme-ports/internal/targets/neovim"
	"bearded-theme-ports/internal/targets/tmtheme"
	"bearded-theme-ports/internal/targets/wezterm"
)

type builderFunc func(root string, inputs buildInputs) ([]string, error)

type buildInputs struct {
	VSCodeThemes []model.ThemeFile
	ZedThemes    []model.ZedThemeFile
}

type targetDefinition struct {
	builder builderFunc
	source  string
}

type buildOptions struct {
	install bool
	targets []string
}

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
	options, err := parseBuildOptions(args)
	if err != nil {
		return err
	}

	inputs, sourcePaths, err := loadInputs(root, options.targets)
	if err != nil {
		return err
	}

	for _, target := range options.targets {
		if _, err := targetsByName[target].builder(root, inputs); err != nil {
			return err
		}
	}

	installedTargets := make([]string, 0, len(options.targets))
	installedPaths := make([]string, 0, len(options.targets))
	if options.install {
		for _, target := range options.targets {
			if !install.SupportedTarget(target) {
				continue
			}
			installPath, err := install.Install(root, target)
			if err != nil {
				return err
			}
			installedTargets = append(installedTargets, target)
			installedPaths = append(installedPaths, installPath)
		}
	}

	commitSHA, err := source.UpstreamCommitSHA(root)
	if err != nil {
		return err
	}

	if err := output.WriteMetadata(root, output.Metadata{
		GeneratedTargets: options.targets,
		InputThemeCount:  inputThemeCount(inputs),
		SourcePath:       primarySourcePath(sourcePaths),
		SourcePaths:      sourcePaths,
		UpstreamCommit:   commitSHA,
		UpstreamRepoURL:  source.UpstreamRepoURL,
	}); err != nil {
		return err
	}

	fmt.Printf("✅ Generated %d themes for targets: %s\n", inputThemeCount(inputs), formatTargets(options.targets))
	if len(installedTargets) > 0 {
		fmt.Printf("📦 Installed generated files for targets: %s\n", formatTargets(installedTargets))
		for index := range installedTargets {
			fmt.Printf("   %s -> %s\n", installedTargets[index], installedPaths[index])
		}
	}
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

func parseBuildOptions(args []string) (buildOptions, error) {
	options := buildOptions{}
	targetArgs := make([]string, 0, len(args))

	for _, arg := range args {
		switch arg {
		case "--install", "-i":
			options.install = true
		default:
			targetArgs = append(targetArgs, arg)
		}
	}

	targets, err := parseTargets(targetArgs)
	if err != nil {
		return buildOptions{}, err
	}
	options.targets = targets
	return options, nil
}

func loadInputs(root string, targets []string) (buildInputs, []string, error) {
	inputs := buildInputs{}
	sourceSet := make(map[string]bool, 2)
	sourcePaths := make([]string, 0, 2)

	needsVSCode := false
	needsZed := false
	for _, target := range targets {
		switch targetsByName[target].source {
		case "vscode":
			needsVSCode = true
		case "zed":
			needsZed = true
		}
	}

	if needsVSCode {
		themes, err := source.LoadThemes(root)
		if err != nil {
			return buildInputs{}, nil, err
		}
		inputs.VSCodeThemes = themes
		if !sourceSet[source.InputThemesDir(root)] {
			sourceSet[source.InputThemesDir(root)] = true
			sourcePaths = append(sourcePaths, source.InputThemesDir(root))
		}
	}

	if needsZed {
		themes, err := source.LoadZedThemes(root)
		if err != nil {
			return buildInputs{}, nil, err
		}
		inputs.ZedThemes = themes
		if !sourceSet[source.ZedThemesPath(root)] {
			sourceSet[source.ZedThemesPath(root)] = true
			sourcePaths = append(sourcePaths, source.ZedThemesPath(root))
		}
	}

	sort.Strings(sourcePaths)
	return inputs, sourcePaths, nil
}

func inputThemeCount(inputs buildInputs) int {
	if len(inputs.ZedThemes) > 0 {
		return len(inputs.ZedThemes)
	}
	return len(inputs.VSCodeThemes)
}

func primarySourcePath(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func parseTargets(args []string) ([]string, error) {
	if len(args) == 0 {
		return allTargets(), nil
	}

	seen := make(map[string]bool, len(args))
	targets := make([]string, 0, len(args))
	for _, target := range args {
		if _, ok := targetsByName[target]; !ok {
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
	targets := make([]string, 0, len(targetsByName))
	for target := range targetsByName {
		targets = append(targets, target)
	}
	sort.Strings(targets)
	return targets
}

var targetsByName = map[string]targetDefinition{
	"helix": {
		source:  "zed",
		builder: func(root string, inputs buildInputs) ([]string, error) { return helix.Build(root, inputs.ZedThemes) },
	},
	"neovim": {
		source:  "zed",
		builder: func(root string, inputs buildInputs) ([]string, error) { return neovim.Build(root, inputs.ZedThemes) },
	},
	"tmtheme": {
		source: "vscode",
		builder: func(root string, inputs buildInputs) ([]string, error) {
			return tmtheme.Build(root, inputs.VSCodeThemes)
		},
	},
	"wezterm": {
		source: "vscode",
		builder: func(root string, inputs buildInputs) ([]string, error) {
			return wezterm.Build(root, inputs.VSCodeThemes)
		},
	},
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
	issues := make([]string, 0, 5)

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
		issues = append(issues, fmt.Sprintf("prepared VS Code theme JSON missing at %s", source.InputThemesDir(root)))
	}

	if _, err := os.Stat(source.ZedThemesPath(root)); err != nil {
		issues = append(issues, fmt.Sprintf("prepared Zed theme JSON missing at %s", source.ZedThemesPath(root)))
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
  bearded-theme-ports prepare-and-build [--install] [targets...]
  bearded-theme-ports build-all [--install] [targets...]
  bearded-theme-ports build [--install] [targets...]
  bearded-theme-ports list themes
  bearded-theme-ports list targets
  bearded-theme-ports doctor`

func printUsage() {
	fmt.Println(usageText)
}
