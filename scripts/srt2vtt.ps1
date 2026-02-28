# Get all .srt files in the current folder
$files = Get-ChildItem *.srt

foreach ($file in $files) {
    # Search for the S01E01 style pattern in the filename
    # \d+ means "one or more digits"
    if ($file.BaseName -match 'S\d+E\d+') {
        
        // $matches[0] contains the text that matched the pattern
        $episodeCode = $matches[0]
        $outputName = $episodeCode + ".vtt"
        
        Write-Host "Converting: $($file.Name) -> $outputName" -ForegroundColor Cyan
        
        # Run ffmpeg -y (overwrite if exists)
        ffmpeg -i "$($file.FullName)" "$outputName" -y
    } 
    else {
        Write-Host "Skipping: $($file.Name) (No S##E## pattern found)" -ForegroundColor Yellow
    }
}

Write-Host "Done! Folder is now clean." -ForegroundColor Green