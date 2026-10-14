$ErrorActionPreference = "Stop"

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports-tmtheme.zip"
$TargetRoot = if ($env:PI_CODING_AGENT_DIR) { $env:PI_CODING_AGENT_DIR } else { Join-Path $HOME ".pi/agent" }
$TargetDir = Join-Path $TargetRoot "extensions/pigment/themes"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("bearded-theme-ports-pi-pigment-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir "bearded-theme-ports-tmtheme.zip"
$ExtractDir = Join-Path $TempDir "extract"

try {
  New-Item -ItemType Directory -Path $TempDir | Out-Null
  New-Item -ItemType Directory -Path $ExtractDir | Out-Null
  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null

  Write-Host "Downloading latest release from $AssetUrl"
  Invoke-WebRequest -Uri $AssetUrl -OutFile $ArchivePath

  Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force
  Copy-Item -Path (Join-Path $ExtractDir "*") -Destination $TargetDir -Recurse -Force

  Write-Host "Installed Pigment TextMate themes into $TargetDir"
  Write-Host "In Pi, run /pigment convert bearded-theme-monokai-stone, then /reload."
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
