# Test filename
$filename = "Fargo (2014) - S03E05 - NEw anme muthgafucka (1080p AMZN WEB-DL x265 Silence).mkv"

# Markers used for extraction
$startMarker = "S03"
$endMarker   = "(10"

# Build regex
$regex = [regex]::Escape($startMarker) + "(.*?)(?=" + [regex]::Escape($endMarker) + ")"

# Run match
if ($filename -match $regex) {
    Write-Host "Matched segment:"
    Write-Host $matches[0]
}
else {
    Write-Host "No match found."
}
