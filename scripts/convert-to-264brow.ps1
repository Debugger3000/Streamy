
# Convert english .mkvs to x264, 8 bit, aac encoding for normal browser viewing... still maintains resolution


# Set your global markers (used for ALL files)
$startMarker = Read-Host "Enter the START marker - Inclusive (e.g., S03)"
$endMarker   = Read-Host "Enter the END marker - Non-inclusive(e.g., (10)"

# Regex: match everything from START up to (but not including) END
$regex = [regex]::Escape($startMarker) + "(.*?)(?=" + [regex]::Escape($endMarker) + ")"


# Ask for output directory
$outputDir = Read-Host "Enter the output directory (e.g., X:\streamy\media\shows\)"
if (!(Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Get all .mkv files in current directory
$files = Get-ChildItem -Path . -Filter *.mkv -File

foreach ($file in $files) {
    # Match pattern like S01E04 (one letter, two digits, one letter, two digits)
    if ($file.Name -match $regex) {
        $episodeCode = $matches[0]

        # Set output filename
        $outputFile = Join-Path $outputDir "$episodeCode.mp4"

        Write-Host "Converting '$($file.Name)' to '$outputFile'..."

        # FFmpeg command as argument array
        $ffmpegArgs = @(
            "-i", $file.FullName,
            "-c:v", "libx264",
            "-pix_fmt", "yuv420p",
            "-c:a", "aac",
            "-b:a", "128k",
            "-preset", "fast",
            $outputFile
        )

        # Run FFmpeg and show its output in this same PowerShell window
        & ffmpeg @ffmpegArgs
    }
    else {
        Write-Host "Skipping '$($file.Name)': filename does not match episode pattern."
    }
}

Write-Host "All done!"
