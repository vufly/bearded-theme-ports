$ErrorActionPreference = "Stop"

$Repo = "vufly/bearded-theme-ports"
$AssetUrl = "https://github.com/$Repo/releases/latest/download/bearded-theme-ports-tmtheme.zip"
$TempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("bearded-theme-ports-bat-" + [System.Guid]::NewGuid().ToString("N"))
$ArchivePath = Join-Path $TempDir "bearded-theme-ports-tmtheme.zip"
$ExtractDir = Join-Path $TempDir "extract"

function Resolve-BatCommand {
  foreach ($name in @("bat", "batcat")) {
    $command = Get-Command $name -ErrorAction SilentlyContinue
    if ($null -ne $command) {
      return $command.Source
    }
  }

  throw "Missing bat executable: need bat or batcat"
}

try {
  $BatBin = Resolve-BatCommand
  $BatConfigDir = & $BatBin --config-dir
  $TargetDir = Join-Path $BatConfigDir "themes"

  New-Item -ItemType Directory -Path $TempDir | Out-Null
  New-Item -ItemType Directory -Path $ExtractDir | Out-Null
  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null

  Write-Host "Downloading latest release from $AssetUrl"
  Invoke-WebRequest -Uri $AssetUrl -OutFile $ArchivePath

  Expand-Archive -Path $ArchivePath -DestinationPath $ExtractDir -Force
  Copy-Item -Path (Join-Path $ExtractDir "*") -Destination $TargetDir -Recurse -Force

  & $BatBin cache --build

  Write-Host "Installed bat themes into $TargetDir"
  Write-Host "Run '$BatBin --list-themes' to inspect available themes"
}
finally {
  if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
  }
}
