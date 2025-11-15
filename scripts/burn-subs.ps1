# grab .ass from all .mkvs, and convert all .mkvs into .mp4s with burned in .ass subtitles
# Creates folder for new .mp4s, and uses basename of .mkv files...

# Create output folder if it doesn't exist
$outputFolder = ".\MP4_Burned"
if (-not (Test-Path $outputFolder)) {
    New-Item -ItemType Directory -Path $outputFolder | Out-Null
}

# Get all MKV files in current directory
Get-ChildItem -Path . -Filter *.mkv | ForEach-Object {

    $mkvFile = $_.FullName
    $baseName = $_.BaseName
    $assFile = "$baseName.english.ass"
    $mp4File = Join-Path $outputFolder "$baseName.mp4"

    Write-Host "Processing $mkvFile ..."

    # Step 1: Extract English ASS subtitles
    ffmpeg -y -i "$mkvFile" -map 0:s:m:language:eng "$assFile"

    # Step 2: Burn subtitles into MP4, output to folder
    ffmpeg -y -i "$mkvFile" -vf "ass=$assFile" -c:v libx264 -crf 23 -preset fast -c:a copy "$mp4File"

    Write-Host "Finished: $mp4File`n"
}