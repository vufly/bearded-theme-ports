$ErrorActionPreference = "Stop"

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports-zellij.zip"
# Install under $HOME\.config\zellij\themes on every platform (honoring
# XDG_CONFIG_HOME when set) so the path stays consistent with darwin and
# linux installs.
if ($env:XDG_CONFIG_HOME) {
  $ConfigRoot = $env:XDG_CONFIG_HOME
} else {
  $ConfigRoot = Join-Path $HOME ".config"
}
$TargetDir = Join-Path $ConfigRoot "zellij\themes"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("bearded-theme-ports-zellij-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir "bearded-theme-ports-zellij.zip"
$ExtractDir = Join-Path $TempDir "extract"

try {
  New-Item -ItemType Directory -Path $TempDir | Out-Null
  New-Item -ItemType Directory -Path $ExtractDir | Out-Null
  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null

  Write-Host "Downloading latest release from $AssetUrl"
  Invoke-WebRequest -Uri $AssetUrl -OutFile $ArchivePath

  Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force
  Copy-Item -Path (Join-Path $ExtractDir "*") -Destination $TargetDir -Recurse -Force

  Write-Host "Installed Zellij themes into $TargetDir"
  Write-Host "Activate one in your zellij config.kdl with:  theme `"<slug>`""
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
