
# Ask for output directory
$outputDir = Read-Host "Enter the output directory (e.g., X:\streamy\media\shows\)"
if (!(Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Get all .mkv files in current directory
$files = Get-ChildItem -Path . -Filter *.mkv -File

foreach ($file in $files) {
    # Match pattern like S01E04 (one letter, two digits, one letter, two digits)
    if ($file.Name -match "[A-Z]\d{2}[A-Z]\d{2}") {
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
