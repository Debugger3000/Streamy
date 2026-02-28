
# Convert english .mkvs to x264, 8 bit, aac encoding for normal browser viewing... still maintains resolution


#si=0 (Stream 0:3): Signs & Songs (Mostly intro/outro text).
#si=1 (Stream 0:4): Dialogue (Non-Honorifics).
#si=2 (Stream 0:5): Dialogue (Honorifics) — This is what you want.

# audio Notes:
	# EAC3 - audio has 6 channels
	# If i try to conver to AAC (2 channels) it silently fails.
	# <"-ac", "2"> - to convert to two channels if need be	


# Set your global markers (used for ALL files)
# $startMarker = Read-Host "Enter the START marker - Inclusive (e.g., S03)"
# $escapedStart = [regex]::Escape($startMarker)
# $endMarkerPattern   = Read-Host "Enter the END marker or press enter for no end marker - Non-inclusive(e.g., (10)"
# $endMarker = ""
# if ([string]::IsNullOrWhiteSpace($endMarkerPattern)) {
#     # Match from start to the very end of the string ($)
#     $endMarker = "$escapedStart(.*)$"
# } 
# else{
#     $endMarker = [regex]::Escape($endMarkerPattern)
# }

# $escapedEnd = [regex]::Escape($endMarker)

# Set your global markers
$startMarker = Read-Host "Enter the START marker - Inclusive (e.g., S01)"
$endMarkerPattern = Read-Host "Enter the END marker or press enter for no end marker"

# 1. Escape the start marker first so it's ready to use
$escapedStart = [regex]::Escape($startMarker)

$finalRegex = ""

if ([string]::IsNullOrWhiteSpace($endMarkerPattern)) {
    # Match from start to the very end of the string ($)
    $finalRegex = "$escapedStart(.*)$"
} 
else {
    $escapedEnd = [regex]::Escape($endMarkerPattern)
    $finalRegex = "$escapedStart(.*?)(?=$escapedEnd)"
}

# Write-Output "Endmarker is: '$($endMarker)'

# Regex: match everything from START up to (but not including) END
#$regex = [regex]::Escape($escapedStart) + "(.*?)(?=" + [regex]::Escape($endMarker) + ")"

# Ask for output directory
$outputDir = Read-Host "Enter the output directory (e.g., X:\streamy\media\shows\)"
if (!(Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}



# # Get all .mkv files in current directory
$files = Get-ChildItem -Path . -Filter *.mkv -File

foreach ($file in $files) {
    if ($file.Name -match $finalRegex) {
        $episodeCode = $matches[0]
        $outputFile = Join-Path $outputDir "$episodeCode.mp4"

        Write-Host "New file name is: '$episodeCode'..."
        Write-Host "Converting '$($file.Name)' to '$outputFile' with Jap Audio + Eng Subs..."

        # Format the file path specifically for the FFmpeg subtitle filter
        # It needs forward slashes and the colon escaped: C\: /Path/To/File.mkv
        $subPath = $file.FullName.Replace("\", "/").Replace(":", "\:")

        $ffmpegArgs = @(
            "-i", $file.FullName,
            "-map", "0:0",                 # Take Video
            "-map", "0:1",                 # Take Japanese Audio (Stream #0:2)
            "-vf", "subtitles='$subPath':si=0",    
# Burn in Stream #0:4 (full dialogue) (defaults to first sub stream) ( guide to whatever one we want) ex: `si=2` = second subtitle track
            "-c:v", "libx264",
            "-pix_fmt", "yuv420p",
            "-c:a", "aac",
            "-b:a", "128k",
            "-ac", "2",
            "-f", "mp4",
            "-preset", "fast",
            "-y",                            # Overwrite output file if it exists
            $outputFile
        )

        & ffmpeg @ffmpegArgs
    }
    else {
        Write-Host "Skipping '$($file.Name)': filename does not match pattern."
    }
}

Write-Host "All done!"
