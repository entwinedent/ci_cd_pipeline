param(
    [string]$Path = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

function Invoke-CargoClean {
    param(
        [string]$ProjectPath
    )

    if ($DryRun) {
        Write-Host "Would run: cargo clean in $ProjectPath"
        return
    }

    Push-Location $ProjectPath
    try {
        & cargo clean
        if ($LASTEXITCODE -ne 0) {
            throw "cargo clean failed in $ProjectPath"
        }
        Write-Host "Cleaned Rust project at $ProjectPath"
    }
    finally {
        Pop-Location
    }
}

$repoRoot = (Resolve-Path $Path).Path
$projects = Get-ChildItem -Path $repoRoot -Recurse -File -Filter "Cargo.toml" -ErrorAction SilentlyContinue |
    Select-Object -ExpandProperty DirectoryName -Unique

if (-not $projects) {
    Write-Host "No Rust projects found under $repoRoot"
    return
}

$cargoCommand = Get-Command cargo -ErrorAction SilentlyContinue
if ($null -eq $cargoCommand) {
    Write-Warning "cargo was not found on PATH. Removing target directories manually instead."
}

foreach ($project in $projects) {
    if ($null -ne $cargoCommand) {
        Invoke-CargoClean -ProjectPath $project
    }
    else {
        $targetPath = Join-Path $project "target"
        if (Test-Path $targetPath) {
            if ($DryRun) {
                Write-Host "Would remove $targetPath"
            }
            else {
                Remove-Item -Path $targetPath -Recurse -Force
                Write-Host "Removed $targetPath"
            }
        }
    }
}

Write-Host "Rust cleanup completed."
