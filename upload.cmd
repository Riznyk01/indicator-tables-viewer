@echo off

set /p FTP_SERVER=<ftp_server.txt
set /p FTP_USERNAME=<ftp_username.txt
set /p FTP_PASSWORD_HASH=<ftp_password_hash.txt

set LOCAL_DIR=build

cd /D %LOCAL_DIR%

curl --ftp-ssl --ftp-ssl-reqd --insecure --user %FTP_USERNAME%:%FTP_PASSWORD_HASH% --upload-file update.zip ftp://%FTP_SERVER%/public_html/viewer_updates/update.zip
curl --ftp-ssl --ftp-ssl-reqd --insecure --user %FTP_USERNAME%:%FTP_PASSWORD_HASH% --upload-file ver ftp://%FTP_SERVER%/public_html/viewer_updates/ver
