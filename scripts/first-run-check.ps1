param(
    [string]$AppPath = "",
    [switch]$NoLaunch
)

$ErrorActionPreference = "Stop"

function Write-Section {
    param([string]$Text)
    Write-Host ""
    Write-Host "== $Text ==" -ForegroundColor Cyan
}

function Test-WebView2Runtime {
    $runtimeKey = "{F3017226-FE2A-4295-8BDF-00C3A9A7E4C5}"
    $paths = @(
        "HKLM:\SOFTWARE\Microsoft\EdgeUpdate\Clients\$runtimeKey",
        "HKLM:\SOFTWARE\WOW6432Node\Microsoft\EdgeUpdate\Clients\$runtimeKey"
    )

    foreach ($path in $paths) {
        if (Test-Path $path) {
            $version = (Get-ItemProperty -Path $path -ErrorAction SilentlyContinue).pv
            return @{ Installed = $true; Version = $version; RegistryPath = $path }
        }
    }

    return @{ Installed = $false; Version = ""; RegistryPath = "" }
}

Write-Host "Indus Personal Work Track - First Run Check" -ForegroundColor Green
Write-Host "Date: $(Get-Date)"

Write-Section "Runtime Checks"
$wv2 = Test-WebView2Runtime
if ($wv2.Installed) {
    Write-Host "WebView2 Runtime: OK ($($wv2.Version))"
} else {
    Write-Host "WebView2 Runtime: MISSING" -ForegroundColor Yellow
    Write-Host "Install from: https://developer.microsoft.com/en-us/microsoft-edge/webview2/"
}

Write-Section "Data Directory Checks"
$appData = if ($env:APPDATA) { $env:APPDATA } else { "." }
$dataDir = Join-Path $appData "indus-task"
$dbPath = Join-Path $dataDir "indus-task.db"

if (-not (Test-Path $dataDir)) {
    New-Item -Path $dataDir -ItemType Directory -Force | Out-Null
    Write-Host "Created data directory: $dataDir"
} else {
    Write-Host "Data directory exists: $dataDir"
}

try {
    $probeFile = Join-Path $dataDir ".write-test"
    "ok" | Set-Content -Path $probeFile -Encoding UTF8
    Remove-Item $probeFile -Force
    Write-Host "Write access: OK"
} catch {
    Write-Host "Write access: FAILED ($($_.Exception.Message))" -ForegroundColor Red
    exit 1
}

if (Test-Path $dbPath) {
    Write-Host "Database file: PRESENT ($dbPath)"
} else {
    Write-Host "Database file: NOT YET CREATED (will be created on app start)"
}

Write-Section "Application Binary Check"
if (-not $AppPath) {
    $scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
    $candidate1 = Join-Path $scriptRoot "..\backend\build\bin\IndusTaskManager.exe"
    $candidate2 = Join-Path $scriptRoot "IndusTaskManager.exe"
    if (Test-Path $candidate1) {
        $AppPath = (Resolve-Path $candidate1).Path
    } elseif (Test-Path $candidate2) {
        $AppPath = (Resolve-Path $candidate2).Path
    }
}

if (-not $AppPath -or -not (Test-Path $AppPath)) {
    Write-Host "App executable not found. Pass -AppPath <path-to-IndusTaskManager.exe>." -ForegroundColor Red
    exit 1
}

$appInfo = Get-Item $AppPath
Write-Host "App executable: $($appInfo.FullName)"
Write-Host "Last modified : $($appInfo.LastWriteTime)"
Write-Host "Size          : $([math]::Round($appInfo.Length / 1MB, 2)) MB"

if (-not $NoLaunch) {
    Write-Section "Launching Application"
    Start-Process -FilePath $AppPath | Out-Null
    Write-Host "Application started successfully."
}

Write-Host ""
Write-Host "First-run check completed." -ForegroundColor Green
