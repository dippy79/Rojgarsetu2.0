#!/usr/bin/env pwsh
<#
.SYNOPSIS
    RojgarSetu2.0 Project Cleanup - Removes junk, reorganizes to services/, moves PHASE TODOs
.PARAMETER Preview
    Dry-run mode (WhatIf)
.PARAMETER Force
    Skip confirmations
#>

param(
    [switch]$Preview,
    [switch]$Force
)

$root = 'f:/Rojgarsetu2.0/rojgarsetu2'
$ErrorActionPreference = 'Stop'

Write-Host "=== RojgarSetu 2.0 CLEANUP ===" -ForegroundColor Green
Write-Host "Root: $root" -ForegroundColor Yellow

# Ensure docs/archive exists
$docsArchive = Join-Path $root 'docs/archive'
if (!(Test-Path $docsArchive)) { New-Item -Path $docsArchive -ItemType Directory -Force | Out-Null }

# Step 1: Move PHASE*_TODO.md to docs/archive/
$phaseFiles = Get-ChildItem -Path $root -Filter 'PHASE*_TODO.md' 
foreach ($file in $phaseFiles) {
    $dest = Join-Path $docsArchive $file.Name
    Write-Host "Move: $($file.Name) -> docs/archive/" -ForegroundColor Cyan
    Move-Item -Path $file.FullName -Destination $dest -Force:$Preview.IsPresent -WhatIf:$Preview
}

# Step 2: Root junk files/folders
$rootJunk = @(
    '$null',
    '1KB', 
    '1KB)"',
    '-Force.bak',
    'backend-backend-fix.bak',
    'Backups created',
    'Copy-Item.bak',
    "backend_go/'\''yyyyMMdd_HHmmss')/''",
    'echo'
)
foreach ($junk in $rootJunk) {
    $path = Join-Path $root $junk
    if (Test-Path $path) {
        Write-Host "Delete root junk: $junk" -ForegroundColor Red
        Remove-Item -Path $path -Recurse -Force -WhatIf:$Preview
    }
}

# Step 3: Backup folders
$backupFolders = @(
    'deployment_backup_*',
    'backend_go/deployment/backup_temp_stubs_*'
)
foreach ($pattern in $backupFolders) {
    $matches = Get-ChildItem -Path $root -Filter $pattern -Directory
    foreach ($match in $matches) {
        Write-Host "Delete backup: $($match.Name)" -ForegroundColor Red
        Remove-Item -Path $match.FullName -Recurse -Force -WhatIf:$Preview
    }
}

# Step 4: .bak, .new, .original files (safe patterns)
$bakPatterns = @(
    'deployment/Dockerfile.backend.*bak*',
    'deployment/Dockerfile.crawler.*bak*',
    'deployment/Dockerfile.backend.new',
    'deployment/Dockerfile.crawler.new',
    'deployment/docker-compose.yml.bak',
    "backend_go/internal/db/*.go.bak",
    "backend_go/cmd/server/main.go.new",
    "auth_service_java/target/auth-service-*.jar.original"
)
foreach ($pattern in $bakPatterns) {
    $matches = Get-ChildItem -Path $root -Filter $pattern -Recurse
    foreach ($match in $matches) {
        Write-Host "Delete .bak/.new: $($match.FullName)" -ForegroundColor Red
        Remove-Item -Path $match.FullName -Force -WhatIf:$Preview
    }
}

# Step 5: Maven target & Go binaries
Write-Host "Delete auth_service_java/target/" -ForegroundColor Red
Remove-Item -Path (Join-Path $root 'auth_service_java/target') -Recurse -Force -WhatIf:$Preview

$goBinaries = @('crawler_go/main.exe')
foreach ($bin in $goBinaries) {
    $path = Join-Path $root $bin
    if (Test-Path $path) {
        Write-Host "Delete binary: $bin" -ForegroundColor Red
        Remove-Item -Path $path -Force -WhatIf:$Preview
    }
}

# Step 6: Temp logs (pattern match incomplete timestamps)
$logJunk = Get-ChildItem -Path (Join-Path $root 'deployment/logs') -Filter '*$(date*' -ErrorAction SilentlyContinue
foreach ($log in $logJunk) {
    Write-Host "Delete temp log: $($log.Name)" -ForegroundColor Red
    Remove-Item -Path $log.FullName -Recurse -Force -WhatIf:$Preview
}

# Step 7: Create services/ & Rename service folders
$serviceRenames = @{
    'backend_go' = 'services/backend-go'
    'crawler_go' = 'services/crawler-go' 
    'api_gateway_node' = 'services/api-gateway-node'
    'auth_service_java' = 'services/auth-java'
    'ai_engine_python' = 'services/ai-engine-python'
}

# Non-destructive mode: do NOT rename folders automatically.
# Instead, ensure `services/` exists and print suggested mappings for manual review.
$servicesDir = Join-Path $root 'services'
if (!(Test-Path $servicesDir)) { New-Item -Path $servicesDir -ItemType Directory -Force | Out-Null }
foreach ($old in $serviceRenames.Keys) {
    $oldPath = Join-Path $root $old
    $newPath = Join-Path $root $serviceRenames[$old]
    if (Test-Path $oldPath -PathType Container) {
        Write-Host "Found: $oldPath -> Suggested location: $newPath (NO ACTION taken)" -ForegroundColor Yellow
    } else {
        Write-Host "Not found: $oldPath (skipping)" -ForegroundColor DarkGray
    }
}

# Step 8: Cleanup empty deployment/logs temp files
$tempFiles = @('deployment/images.txt', 'deployment/build_output.txt', 'deployment/logs/psql_version.txt')
foreach ($file in $tempFiles) {
    $path = Join-Path $root $file
    if (Test-Path $path) {
        Write-Host "Delete temp: $file" -ForegroundColor Red
        Remove-Item -Path $path -Force -WhatIf:$Preview
    }
}

Write-Host "`n=== CLEANUP COMPLETE ===" -ForegroundColor Green
if ($Preview) {
    Write-Host "Run: .\cleanup_project.ps1 -Force to execute (no -Preview)" -ForegroundColor Yellow
} else {
    Write-Host "Project cleaned! Update docker-compose contexts next." -ForegroundColor Green
}
