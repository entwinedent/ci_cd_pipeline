param(
    [string]$Path = $PSScriptRoot,
    [switch]$DryRun
)

$ErrorActionPreference = "Stop"

$targetPath = (Resolve-Path $Path).Path

Write-Host "Scanning for Python cache artifacts under $targetPath"

$cacheDirs = Get-ChildItem -Path $targetPath -Recurse -Directory -Filter "__pycache__" -ErrorAction SilentlyContinue
$cacheFiles = Get-ChildItem -Path $targetPath -Recurse -File -Include "*.pyc", "*.pyo" -ErrorAction SilentlyContinue

if (-not $cacheDirs -and -not $cacheFiles) {
    Write-Host "No Python cache artifacts found."
    return
}

$removedCount = 0

foreach ($dir in $cacheDirs) {
    if ($DryRun) {
        Write-Host "Would remove directory: $($dir.FullName)"
    }
    else {
        Remove-Item -Path $dir.FullName -Recurse -Force
        Write-Host "Removed directory: $($dir.FullName)"
        $removedCount++
    }
}

foreach ($file in $cacheFiles) {
    if ($DryRun) {
        Write-Host "Would remove file: $($file.FullName)"
    }
    else {
        Remove-Item -Path $file.FullName -Force
        Write-Host "Removed file: $($file.FullName)"
        $removedCount++
    }
}

if ($DryRun) {
    $dirCount = if ($cacheDirs) { $cacheDirs.Count } else { 0 }
    $fileCount = if ($cacheFiles) { $cacheFiles.Count } else { 0 }
    Write-Host "Dry run complete. Found $dirCount cache directories and $fileCount cache files."
}
else {
    Write-Host "Cleanup complete. Removed $removedCount items."
}
