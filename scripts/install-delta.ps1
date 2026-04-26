$ErrorActionPreference = "Stop"

# Installs all generated delta gitconfig files so users can include either the
# consolidated file or only the specific variant files they want.

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports-delta.zip"
$TargetDir = Join-Path $env:USERPROFILE ".config/git"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("bearded-theme-ports-delta-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir "bearded-theme-ports-delta.zip"
$ExtractDir = Join-Path $TempDir "extract"

try {
  New-Item -ItemType Directory -Path $TempDir | Out-Null
  New-Item -ItemType Directory -Path $ExtractDir | Out-Null
  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null

  Write-Host "Downloading latest release from $AssetUrl"
  Invoke-WebRequest -Uri $AssetUrl -OutFile $ArchivePath

  Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force

  $SourceFile = Join-Path $ExtractDir "bearded-theme.gitconfig"
  if (-not (Test-Path $SourceFile)) {
    Write-Error "Delta gitconfig files missing from release asset"
    exit 1
  }

  Copy-Item -Path (Join-Path $ExtractDir "*") -Destination $TargetDir -Recurse -Force
  Write-Host "Installed delta presets into $TargetDir"

  Write-Host ""
  Write-Host "Next steps:"
  Write-Host "  1. Optionally add the manual [include] example from this repo's README."
  Write-Host "  2. Make sure delta is your pager:"
  Write-Host "       git config --global core.pager delta"
  Write-Host "       git config --global interactive.diffFilter `"delta --color-only`""
  Write-Host "  3. Activate a variant by name, for example:"
  Write-Host "       git config --global delta.features bearded-theme-monokai-stone"
  Write-Host ""
  Write-Host "  Installed files include:"
  Write-Host "       $TargetDir\bearded-theme-monokai-stone.gitconfig"
  Write-Host "       $TargetDir\bearded-theme.gitconfig"
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
