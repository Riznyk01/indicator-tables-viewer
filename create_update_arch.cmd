@echo off
set build_folder=build
set update_zip=update.zip

if exist "%~dp0%build_folder%\%update_zip%" (
    del "%~dp0%build_folder%\%update_zip%"
)

cd /D "%build_folder%"
"C:\Program Files\7-Zip\7z" a "../"%build_folder%"/%update_zip%" "viewer.exe" "resources" "config.toml"
