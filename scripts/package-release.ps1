param(
    [string]$Version = ""
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$backendDir = Join-Path $repoRoot "backend"
$frontendDir = Join-Path $backendDir "frontend"
$toolsDir = Join-Path $repoRoot "tools"
$portableNsis = Join-Path $toolsDir "nsis-3.11\Bin"
$portableNsisExe = Join-Path $portableNsis "makensis.exe"
$portableNsisZip = Join-Path $toolsDir "nsis-3.11.zip"

if (-not (Get-Command makensis -ErrorAction SilentlyContinue) -and -not (Test-Path $portableNsisExe)) {
    New-Item -ItemType Directory -Path $toolsDir -Force | Out-Null
    Write-Host "Downloading portable NSIS..." -ForegroundColor Yellow
    curl.exe -L "https://sourceforge.net/projects/nsis/files/NSIS%203/3.11/nsis-3.11.zip/download" -o $portableNsisZip
    Expand-Archive -Path $portableNsisZip -DestinationPath $toolsDir -Force
}

if (Test-Path $portableNsisExe) {
    $env:Path = "$portableNsis;$env:Path"
    Write-Host "Using portable NSIS from $portableNsis"
}

if (-not $Version) {
    $wailsConfig = Get-Content (Join-Path $backendDir "wails.json") -Raw | ConvertFrom-Json
    $Version = $wailsConfig.version
}

Write-Host "Packaging release version $Version" -ForegroundColor Cyan

Push-Location $frontendDir
npm run build
Pop-Location

Push-Location $backendDir
& "$env:USERPROFILE\go\bin\wails.exe" build -clean -nsis
Pop-Location

$releaseRoot = Join-Path $repoRoot "release"
$releaseDir = Join-Path $releaseRoot "indus-personal-work-track-v$Version"
$zipPath = Join-Path $releaseRoot "indus-personal-work-track-v$Version-windows.zip"

if (Test-Path $releaseDir) { Remove-Item $releaseDir -Recurse -Force }
if (Test-Path $zipPath) { Remove-Item $zipPath -Force }
New-Item -Path $releaseDir -ItemType Directory -Force | Out-Null

$exePath = Join-Path $backendDir "build\bin\IndusTaskManager.exe"
if (-not (Test-Path $exePath)) {
    throw "Release executable not found at $exePath"
}
Copy-Item $exePath -Destination (Join-Path $releaseDir "IndusTaskManager.exe") -Force

$installerCandidates = Get-ChildItem (Join-Path $backendDir "build\bin") -Filter "*installer*.exe" -ErrorAction SilentlyContinue
if ($installerCandidates) {
    $installer = $installerCandidates | Sort-Object LastWriteTime -Descending | Select-Object -First 1
    Copy-Item $installer.FullName -Destination (Join-Path $releaseDir $installer.Name) -Force
    Write-Host "Included installer: $($installer.Name)"
} else {
    Write-Host "Installer artifact not found, packaging portable exe only." -ForegroundColor Yellow
}

Copy-Item (Join-Path $repoRoot "scripts\first-run-check.ps1") -Destination (Join-Path $releaseDir "first-run-check.ps1") -Force
Copy-Item (Join-Path $repoRoot "scripts\first-run-check.cmd") -Destination (Join-Path $releaseDir "first-run-check.cmd") -Force

$userGuide = @"
Indus Personal Work Track - Quick Start
======================================

1) Run first-run-check.cmd
2) Open IndusTaskManager.exe
3) Login with:
   - admin
   - developer
   - reporter

Data path:
  %APPDATA%\indus-task\indus-task.db
"@
$userGuide | Set-Content -Path (Join-Path $releaseDir "QUICK-START.txt") -Encoding UTF8

Compress-Archive -Path (Join-Path $releaseDir "*") -DestinationPath $zipPath

Write-Host ""
Write-Host "Release folder: $releaseDir" -ForegroundColor Green
Write-Host "Release zip   : $zipPath" -ForegroundColor Green
