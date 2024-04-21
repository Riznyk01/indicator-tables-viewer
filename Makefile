build:
	go build -ldflags -H=windowsgui -o viewer.exe cmd/viewer/main.go cmd/viewer/theme.go
	update_ver.cmd
buildcli:
	go build -o viewer.exe cmd/viewer/main.go cmd/viewer/theme.go
	update_ver.cmd
buildlauncher:
	go build -ldflags -H=windowsgui -o launcher.exe cmd/launcher/main.go
	update_ver.cmd
buildlaunchercli:
	go build -o launcher.exe cmd/launcher/main.go
	update_ver.cmd
