$ErrorActionPreference = "Stop"

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports-codex.zip"
$TargetRoot = if ($env:CODEX_HOME) { $env:CODEX_HOME } else { Join-Path $HOME ".codex" }
$TargetDir = Join-Path $TargetRoot "themes"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("bearded-theme-ports-codex-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir "bearded-theme-ports-codex.zip"
$ExtractDir = Join-Path $TempDir "extract"

try {
  New-Item -ItemType Directory -Path $TempDir | Out-Null
  New-Item -ItemType Directory -Path $ExtractDir | Out-Null
  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null

  Write-Host "Downloading latest release from $AssetUrl"
  Invoke-WebRequest -Uri $AssetUrl -OutFile $ArchivePath

  Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force
  Copy-Item -Path (Join-Path $ExtractDir "*") -Destination $TargetDir -Recurse -Force

  Write-Host "Installed Codex themes into $TargetDir"
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
