if (Test-Path .env) {
    Get-Content .env | Where-Only { $_ -and -not $_.StartsWith("#") } | ForEach-Object {
        $name, $value = $_.Split('=', 2)
        [System.Environment]::SetEnvironmentVariable($name.Trim(), $value.Trim())
    }
}

if (-not $env:DATABASE_URL) { $env:DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5435/rojgarsetu2?sslmode=disable" }
if (-not $env:REDIS_URL) { $env:REDIS_URL="redis://localhost:6380" }
if (-not $env:JWT_SECRET) { $env:JWT_SECRET="super-secret-jwt-key-minimum-32-characters-long" }
if (-not $env:REFRESH_TOKEN_KEY) { $env:REFRESH_TOKEN_KEY="super-secret-refresh-key-minimum-32-chars" }
if (-not $env:CORS_ORIGINS) { $env:CORS_ORIGINS="http://localhost:8080,http://localhost:3000,http://localhost:3001" }
if (-not $env:PORT) { $env:PORT="8083" }
if (-not $env:MIGRATIONS_PATH) { $env:MIGRATIONS_PATH="file://./migrations" }

go run cmd/server/main.go
