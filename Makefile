.PHONY: build buildlauncher

build:
	go build -ldflags -H=windowsgui -o build/viewer.exe cmd/viewer/main.go cmd/viewer/theme.go cmd/viewer/icon.go && update_ver.cmd
buildlauncher:
	go build -ldflags -H=windowsgui -o build/launcher.exe cmd/launcher/main.go && update_ver.cmd