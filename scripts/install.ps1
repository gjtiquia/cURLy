#!/usr/bin/env pwsh
# cURLy installer for Windows (PowerShell)
# Converted from scripts/install.sh

$ErrorActionPreference = "Stop"

Write-Host "Installing cURLy! The snake game with cURL"
Write-Host ""

$REPO = "gjtiquia/cURLy"
$GITHUB_API = "https://api.github.com/repos/$REPO"

# Detect architecture via CIM (more reliable than PROCESSOR_ARCHITECTURE)
$SystemType = (Get-CimInstance Win32_ComputerSystem).SystemType
switch -Regex ($SystemType) {
    "X86-based"   { $ARCH = "386" }
    "x64-based"   { $ARCH = "amd64" }
    "ARM64-based" { $ARCH = "arm64" }
    default       {
        Write-Host "Install Failed: cURLy for Windows is only available for 386, amd64, or ARM64. Detected: $SystemType" -ForegroundColor Red
        return 1
    }
}

$OS = "windows"
$PLATFORM = "$OS/$ARCH"
Write-Host "Detected platform: $PLATFORM"

# Get latest release
try {
    $release = Invoke-RestMethod -Uri "$GITHUB_API/releases/latest" -Method Get
    $VERSION = $release.tag_name
} catch {
    Write-Host "Install Failed: Could not fetch latest release information." -ForegroundColor Red
    return 1
}

if ([string]::IsNullOrEmpty($VERSION)) {
    Write-Host "Install Failed: Could not determine latest version." -ForegroundColor Red
    return 1
}

$VERSION_NO_V = $VERSION.TrimStart('v')
Write-Host "Version: $VERSION_NO_V" -ForegroundColor Blue

# Construct download URL
$ASSET_NAME = "cURLy_${OS}_${ARCH}.exe"
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/$VERSION/$ASSET_NAME"

# Download
Write-Host ""
Write-Host "Downloading $ASSET_NAME from $DOWNLOAD_URL..."
try {
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $ASSET_NAME -UseBasicParsing
} catch {
    Write-Host "Install Failed: Could not download $DOWNLOAD_URL" -ForegroundColor Red
    Write-Host "Platform $PLATFORM may be unsupported, or check your network." -ForegroundColor Red
    return 1
}

if (-not (Test-Path $ASSET_NAME)) {
    Write-Host "Install Failed: Downloaded file is missing. Did an antivirus remove it?" -ForegroundColor Red
    return 1
}

Write-Host "Download successful!"

# Rename downloaded file
$EXEC_NAME = "cURLy.exe"
Move-Item -Path $ASSET_NAME -Destination $EXEC_NAME -Force

Write-Host ""
Write-Host "Run the following command:"
Write-Host ""
Write-Host ".\$EXEC_NAME" -ForegroundColor Green
Write-Host ""
