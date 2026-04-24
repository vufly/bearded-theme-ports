$ErrorActionPreference = "Stop"

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports.zip"
$TargetDir = Join-Path $HOME ".config/wezterm/themes/bearded-theme"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("bearded-theme-ports-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir "bearded-theme-ports.zip"
$ExtractDir = Join-Path $TempDir "extract"

try {
  New-Item -ItemType Directory -Path $TempDir | Out-Null
  New-Item -ItemType Directory -Path $ExtractDir | Out-Null
  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null

  Write-Host "Downloading latest release from $AssetUrl"
  Invoke-WebRequest -Uri $AssetUrl -OutFile $ArchivePath

  Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force

  $WezTermDir = Join-Path $ExtractDir "wezterm"
  if (-not (Test-Path $WezTermDir)) {
    throw "WezTerm themes were not found in the archive"
  }

  Copy-Item -Path (Join-Path $WezTermDir "*") -Destination $TargetDir -Recurse -Force
  Write-Host "Installed WezTerm themes into $TargetDir"
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
