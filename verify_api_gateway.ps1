$ErrorActionPreference = 'Stop'
Set-Location 'f:\Rojgarsetu2.0\rojgarsetu2\services\api-gateway-node'
$env:PORT = '8080'
$env:CSRF_SECRET = 'dummysecret'
$env:REDIS_URL = ''
$proc = Start-Process -FilePath 'node' -ArgumentList 'src/index.js' -WorkingDirectory (Get-Location) -PassThru
Start-Sleep -Seconds 3
try {
    $get = Invoke-WebRequest -Uri 'http://localhost:8080/api/csrf-token' -Method GET -TimeoutSec 20
    Write-Host "GET_STATUS=$($get.StatusCode)"
    Write-Host "GET_BODY=$($get.Content | Out-String)"

    try {
        $post = Invoke-WebRequest -Uri 'http://localhost:8080/api/some-route' -Method POST -Headers @{ 'Content-Type' = 'application/json' } -Body '{}' -TimeoutSec 20
        Write-Host "POST_STATUS=$($post.StatusCode)"
    }
    catch {
        $response = $_.Exception.Response
        $status = if ($null -ne $response) { $response.StatusCode.value__ } else { 'NO_RESPONSE' }
        $body = if ($null -ne $response) { [System.IO.StreamReader]::new($response.GetResponseStream()).ReadToEnd() } else { $_.Exception.Message }
        Write-Host "POST_STATUS=$status"
        Write-Host "POST_BODY=$body"
    }
}
finally {
    if ($null -ne $proc -and -not $proc.HasExited) { Stop-Process -Id $proc.Id -Force }
}
