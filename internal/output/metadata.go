package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"bearded-theme-ports/internal/source"
)

type Metadata struct {
	GeneratedTargets []string
	InputThemeCount  int
	SourcePath       string
	UpstreamCommit   string
	UpstreamRepoURL  string
}

type manifestFile struct {
	GeneratedAt      string   `json:"generated_at"`
	GeneratedTargets []string `json:"generated_targets"`
	InputThemeCount  int      `json:"input_theme_count"`
	SourcePath       string   `json:"source_path"`
	UpstreamCommit   string   `json:"upstream_commit"`
}

type upstreamFile struct {
	CommitSHA  string `json:"commit_sha"`
	PreparedAt string `json:"prepared_at"`
	RepoURL    string `json:"repo_url"`
	SourcePath string `json:"source_path"`
}

func WriteMetadata(root string, metadata Metadata) error {
	if err := os.MkdirAll(source.MetadataDir(root), 0o755); err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)

	if err := writeJSON(filepath.Join(source.MetadataDir(root), "manifest.json"), manifestFile{
		GeneratedAt:      now,
		GeneratedTargets: metadata.GeneratedTargets,
		InputThemeCount:  metadata.InputThemeCount,
		SourcePath:       metadata.SourcePath,
		UpstreamCommit:   metadata.UpstreamCommit,
	}); err != nil {
		return err
	}

	return writeJSON(filepath.Join(source.MetadataDir(root), "upstream.json"), upstreamFile{
		CommitSHA:  metadata.UpstreamCommit,
		PreparedAt: now,
		RepoURL:    metadata.UpstreamRepoURL,
		SourcePath: metadata.SourcePath,
	})
}

func writeJSON(path string, value any) error {
	content, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	content = append(content, '\n')
	return os.WriteFile(path, content, 0o644)
}
