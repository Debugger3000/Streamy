# Create subs folder if it doesn't exist
$subsPath = ".\subs"
if (-not (Test-Path $subsPath)) {
    New-Item -ItemType Directory -Path $subsPath | Out-Null
}

# Loop through all .mkv files in current directory
Get-ChildItem -Filter "*.mkv" | ForEach-Object {
    $inputFile = $_.FullName
    $baseName = $_.BaseName
    $outputFile = Join-Path $subsPath ("$baseName.vtt")

    Write-Host "Extracting English subtitles from $inputFile → $outputFile"

    ffmpeg -i "$inputFile" -map 0:m:language:eng "$outputFile"
}