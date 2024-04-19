build:
	go build -ldflags -H=windowsgui -o viewer.exe cmd/viewer/main.go cmd/viewer/theme.go
	update_ver.cmd
buildwithcli:
	go build -o viewer.exe cmd/viewer/main.go cmd/viewer/theme.go
	update_ver.cmd
buildlauncher:
	go build -ldflags -H=windowsgui -o launcher.exe cmd/launcher/main.go
buildlaunchercli:
	go build -o launcher.exe cmd/launcher/main.go
