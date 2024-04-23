.PHONY: build buildcli buildlauncher buildlaunchercli

build:
	go build -ldflags -H=windowsgui -o build/viewer.exe cmd/viewer/main.go cmd/viewer/theme.go && update_ver.cmd
buildcli:
	go build -o build/viewer.exe cmd/viewer/main.go cmd/viewer/theme.go && update_ver.cmd
buildlauncher:
	go build -ldflags -H=windowsgui -o build/launcher.exe cmd/launcher/main.go && update_ver.cmd
buildlaunchercli:
	go build -o build/launcher.exe cmd/launcher/main.go && update_ver.cmd