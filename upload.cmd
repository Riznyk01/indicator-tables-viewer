@echo off

set /p FTP_SERVER=<ftp_server.txt
set /p FTP_USERNAME=<ftp_username.txt
set /p FTP_PASSWORD=<ftp_password.txt

set LOCAL_DIR=build

cd /D %LOCAL_DIR%

echo user %FTP_USERNAME%> ftpcmd.dat
echo %FTP_PASSWORD%>> ftpcmd.dat
echo binary>> ftpcmd.dat
echo cd public_html/viewer_updates>> ftpcmd.dat
echo put update.zip>> ftpcmd.dat
echo put ver>> ftpcmd.dat
echo quit>> ftpcmd.dat

ftp -n -s:ftpcmd.dat %FTP_SERVER%

del ftpcmd.dat
