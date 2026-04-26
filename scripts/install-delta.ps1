$ErrorActionPreference = "Stop"

# Installs the consolidated bearded-theme.gitconfig and registers it via
# `git config --global --add include.path`. After running, set
# `[delta] features = bearded-theme-<slug>` to activate a variant.

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
  Write-Error "Missing git executable"
  exit 1
}

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports-delta.zip"
$TargetDir = Join-Path $env:USERPROFILE ".config/git"
$TargetFile = Join-Path $TargetDir "bearded-theme.gitconfig"
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
    Write-Error "Consolidated gitconfig missing from release asset"
    exit 1
  }

  Copy-Item -Path $SourceFile -Destination $TargetFile -Force
  Write-Host "Installed delta presets into $TargetFile"

  $existing = & git config --global --get-all include.path 2>$null
  if ($existing -and ($existing -split "`n" | Where-Object { $_ -eq $TargetFile })) {
    Write-Host "include.path already set; skipping git config update"
  } else {
    & git config --global --add include.path $TargetFile
    Write-Host "Registered include.path = $TargetFile in your global git config"
  }

  Write-Host ""
  Write-Host "Next steps:"
  Write-Host "  1. Make sure delta is your pager:"
  Write-Host "       git config --global core.pager delta"
  Write-Host "       git config --global interactive.diffFilter `"delta --color-only`""
  Write-Host "  2. Activate a variant by name, for example:"
  Write-Host "       git config --global delta.features bearded-theme-monokai-stone"
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
