Write-Host "JSSON - Quick Instructions" -ForegroundColor Cyan
Write-Host "" 
Write-Host "Build CLI (Windows amd64):" -ForegroundColor Yellow
Write-Host "  $env:GOOS=`"windows`"; $env:GOARCH=`"amd64`"; go build -o cmd\jsson\jsson.exe cmd\jsson\main.go"

Write-Host "Build WASM (for playground):" -ForegroundColor Yellow
Write-Host "  $env:GOOS=`"js`"; $env:GOARCH=`"wasm`"; go build -o apps/www/public/jsson.wasm cmd/wasm/main.go"

Write-Host "Package VSIX (use vsce via npx):" -ForegroundColor Yellow
Write-Host "  npx --yes vsce package"

Write-Host "Done."