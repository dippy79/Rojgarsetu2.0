# Local Postgres setup for RojgarSetu 2.0 - FIXED VERSION
# Run: powershell deployment/setup_local_db.ps1
# Assumes Postgres installed with psql on PATH

param(
  [string]$LogFile = 'logs/setup_local_db.txt',
  [switch]$NoLog
)

$log = Join-Path (Split-Path $PSScriptRoot) $LogFile
if ($NoLog) { $log = $null }
$db_host = 'localhost'
$db_superuser = 'postgres'
$db_user = 'rojgarsetu'
$db_pass = 'Dippy79'
$db_name = 'rojgarsetu'
$schema_file = Join-Path (Split-Path $PSScriptRoot) 'database/schema_v3.sql'

if ($log) {
  $logDir = Split-Path $log
  New-Item -Path $logDir -ItemType Directory -Force | Out-Null
  Write-Output \"[$(Get-Date)] Starting FIXED DB setup\" | Tee-Object -FilePath $log -Append
}

if (!(Get-Command psql -ErrorAction SilentlyContinue)) {
  if ($log) { Write-Output \"ERROR: psql not found. Install Postgres.\" | Tee-Object -FilePath $log -Append }
  exit 1
}

# 1. Create user IF NOT EXISTS using DO block
$user_sql = @"
DO `$\$`
BEGIN
  CREATE USER $db_user PASSWORD '$db_pass';
EXCEPTION
  WHEN duplicate_object THEN RAISE NOTICE 'User $db_user already exists';
END
`$\$`;
"@
if ($log) { Write-Output \"Creating user $db_user ...\" | Tee-Object -FilePath $log -Append }
$output = & psql -U $db_superuser -h $db_host -c \"$user_sql\" 2>&1
if ($log) { $output | Tee-Object -FilePath $log -Append }

# 2. Create DB IF NOT EXISTS using DO block
$db_sql = @"
DO `$\$`
BEGIN
  CREATE DATABASE $db_name OWNER $db_user;
EXCEPTION
  WHEN duplicate_database THEN RAISE NOTICE 'Database $db_name already exists';
END
`$\$`;
"@
if ($log) { Write-Output \"Creating database $db_name ...\" | Tee-Object -FilePath $log -Append }
$output = & psql -U $db_superuser -h $db_host -c \"$db_sql\" 2>&1
if ($log) { $output | Tee-Object -FilePath $log -Append }

# 3. Apply schema (idempotent)
if ($log) { Write-Output \"Applying schema from $schema_file ...\" | Tee-Object -FilePath $log -Append }
$output = & psql -U $db_user -h $db_host -d $db_name -f \"$schema_file\" 2>&1
if ($log) { $output | Tee-Object -FilePath $log -Append }

if ($log) { Write-Output \"SUCCESS: DB '$db_name' ready with schema_v3. Log: $log\" | Tee-Object -FilePath $log -Append }
if ($log) { Write-Output \"Test: psql -U $db_user -h $db_host -d $db_name -c \\\"SELECT 1;\\\" \" | Tee-Object -FilePath $log -Append }

exit 0

