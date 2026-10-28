
param(
    [string]$UbuntuIP = "192.168.56.101",
    [string]$SSHUser = "vboxuser"
)

$ErrorActionPreference = "Stop"

function Test-LocalTcpPort {
    param([int]$Port)

    $client = [System.Net.Sockets.TcpClient]::new()

    try {
        $client.Connect("127.0.0.1", $Port)
        return $true
    }
    catch {
        return $false
    }
    finally {
        $client.Dispose()
    }
}

# Load database configuration from .env.local
$localEnvPath = Join-Path $PSScriptRoot ".env.local"

if (-not (Test-Path -LiteralPath $localEnvPath)) {
    throw "Missing .env.local. Create it in the project root."
}

$databaseUrlLine = Get-Content -LiteralPath $localEnvPath |
    Where-Object { $_ -match '^\s*DATABASE_URL\s*=' } |
    Select-Object -First 1

if (-not $databaseUrlLine) {
    throw "DATABASE_URL is missing from .env.local."
}

$databaseUrl = ($databaseUrlLine -split '=', 2)[1].Trim().Trim('"').Trim("'")

if ([string]::IsNullOrWhiteSpace($databaseUrl)) {
    throw "DATABASE_URL is empty in .env.local."
}

# Check which SSH tunnels are already listening
$forwardings = @()

if (-not (Test-LocalTcpPort -Port 23750)) {
    $forwardings += "127.0.0.1:23750:/var/run/docker.sock"
}

if (-not (Test-LocalTcpPort -Port 15432)) {
    $forwardings += "127.0.0.1:15432:127.0.0.1:5432"
}

# Start missing SSH tunnels
if ($forwardings.Count -gt 0) {

    $sshArgs = @(
        "ssh.exe",
        "-N",
        "-o",
        "ExitOnForwardFailure=yes"
    )

    foreach ($forwarding in $forwardings) {
        $sshArgs += @("-L", $forwarding)
    }

    $sshArgs += "$SSHUser@$UbuntuIP"

    $sshCommand = $sshArgs -join " "

    Start-Process `
        -FilePath "powershell.exe" `
        -ArgumentList @("-NoExit", "-Command", $sshCommand)

    Write-Host "Opened SSH tunnel window."
    Write-Host "Enter your Ubuntu password there."
}
else {
    Write-Host "Both local tunnel ports are already listening."
}

# Configure Testcontainers
$env:DOCKER_HOST = "tcp://127.0.0.1:23750"
$env:TESTCONTAINERS_HOST_OVERRIDE = $UbuntuIP
$env:TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE = "/var/run/docker.sock"

# Configure application database
$env:DATABASE_URL = $databaseUrl

# Switch to project root
Set-Location $PSScriptRoot

Write-Host ""
Write-Host "Development environment configured."
Write-Host "Ubuntu host: $UbuntuIP"
Write-Host "Docker tunnel: 127.0.0.1:23750"
Write-Host "PostgreSQL tunnel: 127.0.0.1:15432"
Write-Host "DATABASE_URL loaded from .env.local"
Write-Host ""
Write-Host "Check Docker:"
Write-Host "curl.exe http://127.0.0.1:23750/_ping"
Write-Host ""
Write-Host "Run application:"
Write-Host "go run ./cmd/application-service"
Write-Host ""
Write-Host "Run tests:"
Write-Host "go test ./... -count=1 -v"
