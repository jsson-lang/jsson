# Build script for JSSON v0.0.4
# Builds binaries for Windows, Linux, and macOS

$VERSION = "v0.0.4"
$OUTPUT_DIR = "dist/$VERSION"

# Create output directory
New-Item -ItemType Directory -Force -Path $OUTPUT_DIR | Out-Null

Write-Host "Building JSSON $VERSION for multiple platforms..." -ForegroundColor Green

# Windows AMD64
Write-Host "`nBuilding for Windows (amd64)..." -ForegroundColor Cyan
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o "$OUTPUT_DIR/jsson-$VERSION-windows-amd64.exe" ./cmd/jsson

# Linux AMD64
Write-Host "Building for Linux (amd64)..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o "$OUTPUT_DIR/jsson-$VERSION-linux-amd64" ./cmd/jsson

# macOS AMD64 (Intel)
Write-Host "Building for macOS (amd64)..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o "$OUTPUT_DIR/jsson-$VERSION-darwin-amd64" ./cmd/jsson

# macOS ARM64 (Apple Silicon)
Write-Host "Building for macOS (arm64)..." -ForegroundColor Cyan
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o "$OUTPUT_DIR/jsson-$VERSION-darwin-arm64" ./cmd/jsson

Write-Host "`n✅ Build complete! Binaries are in the '$OUTPUT_DIR' directory" -ForegroundColor Green
Write-Host "`nFiles created:" -ForegroundColor Yellow
Get-ChildItem $OUTPUT_DIR | ForEach-Object {
    $size = [math]::Round($_.Length / 1MB, 2)
    Write-Host "  - $($_.Name) ($size MB)"
}
