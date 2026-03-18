# Cleanup preview + safe delete script for backend_go
# Save as run_cleanup_report.ps1 and run from repo root or paste into PowerShell

$root = "F:\Rojgarsetu2.0\rojgarsetu2\backend_go"
$report = Join-Path $root "deployment\logs\cleanup_report.txt"
New-Item -Path (Split-Path $report) -ItemType Directory -Force | Out-Null
"`n=== Cleanup report generated at $(Get-Date) ===`n" | Out-File $report -Encoding UTF8

Set-Location $root

# 1) Git status (untracked/temp files)
"`n--- Git status (porcelain) ---`n" | Out-File $report -Append
git status --porcelain 2>&1 | Out-File $report -Append

# 2) Go env and cache sizes
"`n--- Go env and cache sizes ---`n" | Out-File $report -Append
"go env GOMODCACHE: $(go env GOMODCACHE)" | Out-File $report -Append
"go env GOCACHE: $(go env GOCACHE)" | Out-File $report -Append
try {
  (Get-ChildItem -Path (go env GOMODCACHE) -Recurse -ErrorAction Stop | Measure-Object -Property Length -Sum).Sum | ForEach-Object { "GOMODCACHE bytes: $_" } | Out-File $report -Append
} catch { "GOMODCACHE not accessible or empty" | Out-File $report -Append }
try {
  (Get-ChildItem -Path (go env GOCACHE) -Recurse -ErrorAction Stop | Measure-Object -Property Length -Sum).Sum | ForEach-Object { "GOCACHE bytes: $_" } | Out-File $report -Append
} catch { "GOCACHE not accessible or empty" | Out-File $report -Append }

# 3) Project artifacts and temp stubs
"`n--- Project artifacts and temp stubs (preview) ---`n" | Out-File $report -Append
Get-ChildItem -Path $root -Recurse -Include *.exe,*.o,*.a,*.tmp,*.log,*.bak,models_temp_for_build.go,models_extra_for_build.go,db_methods_stubs_for_build.go -ErrorAction SilentlyContinue |
  Select-Object FullName,Length,LastWriteTime |
  Format-Table -AutoSize | Out-String | Out-File $report -Append

# 4) Temporary stubs specifically
"`n--- Temporary stub files in internal\db ---`n" | Out-File $report -Append
Get-ChildItem -Path (Join-Path $root "internal\db") -Include "*temp*.go","*extra*.go","*stubs*.go" -ErrorAction SilentlyContinue |
  Select-Object FullName,Length,LastWriteTime | Format-Table -AutoSize | Out-String | Out-File $report -Append

# 5) Docker preview
"`n--- Docker preview ---`n" | Out-File $report -Append
try {
  docker ps -a --format "table {{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Names}}" | Out-File $report -Append
  docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.Size}}" | Out-File $report -Append
  docker images -f "dangling=true" --format "table {{.Repository}}\t{{.ID}}\t{{.Size}}" | Out-File $report -Append
  docker volume ls -f "dangling=true" --format "table {{.Name}}\t{{.Driver}}" | Out-File $report -Append
} catch {
  "Docker not available or permission denied; skip Docker preview" | Out-File $report -Append
}

# 6) Duplicate files (hash groups) - limited to files >1KB to speed up
"`n--- Duplicate files (SHA256 groups) ---`n" | Out-File $report -Append
$hashes = @{}
Get-ChildItem -Path $root -Recurse -File -ErrorAction SilentlyContinue | Where-Object { $_.Length -gt 1024 } |
  ForEach-Object {
    try {
      $h = (Get-FileHash -Algorithm SHA256 -Path $_.FullName).Hash
      if (-not $hashes.ContainsKey($h)) { $hashes[$h] = @() }
      $hashes[$h] += $_.FullName
    } catch {}
  }
$dups = $hashes.GetEnumerator() | Where-Object { $_.Value.Count -gt 1 }
if ($dups.Count -eq 0) {
  "No duplicates found (files >1KB)" | Out-File $report -Append
} else {
  foreach ($g in $dups) {
    "---- Duplicate group ----" | Out-File $report -Append
    $g.Value | ForEach-Object { $_ | Out-File $report -Append }
  }
}

# 7) Summary counts
"`n--- Summary counts ---`n" | Out-File $report -Append
"Temp stubs found: $((Get-ChildItem -Path (Join-Path $root 'internal\db') -Include '*temp*.go','*extra*.go','*stubs*.go' -ErrorAction SilentlyContinue | Measure-Object).Count)" | Out-File $report -Append
"Build logs found: $((Get-ChildItem -Path (Join-Path $root 'deployment\logs') -Include *.txt -ErrorAction SilentlyContinue | Measure-Object).Count)" | Out-File $report -Append

# Show report to user
Write-Output "`nCleanup preview complete. Report saved to:`n$report`n"
Get-Content $report | Write-Output

# Confirm deletion choices
$confirm = Read-Host "`nDo you want to remove temporary stubs and run go clean and docker system prune? Type 'yes' to proceed, anything else to cancel"
if ($confirm -ne "yes") {
  Write-Output "No changes made. Review the report and run the script again when ready."
  exit 0
}

# Perform safe deletions (only the temporary files we created earlier)
Write-Output "`nRemoving temporary stub files..."
Get-ChildItem -Path (Join-Path $root 'internal\db') -Include "*temp*.go","*extra*.go","*stubs*.go" -ErrorAction SilentlyContinue |
  ForEach-Object {
    Write-Output "Removing $($_.FullName)"
    Remove-Item $_.FullName -Force -ErrorAction SilentlyContinue
  }

# Run go clean commands
Write-Output "`nRunning go clean -modcache and go clean -cache..."
go clean -modcache
go clean -cache

# Docker prune (only if docker available)
try {
  Write-Output "`nRunning docker system prune --all --volumes --force..."
  docker system prune --all --volumes --force
} catch {
  Write-Output "Docker prune skipped or failed (docker may not be available)."
}

Write-Output "`nCleanup actions complete. Re-run build to verify."

