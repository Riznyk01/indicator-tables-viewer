build:
	go build -ldflags -H=windowsgui -o viewer.exe cmd/viewer/main.go cmd/viewer/theme.go

buildwithcli:
	go build -o viewer.exe cmd/viewer/main.go cmd/viewer/theme.go