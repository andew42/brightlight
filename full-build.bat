@echo off
rem Local equivalent of the CI pipeline (.github/workflows/build.yml):
rem builds the frontend and backend the same way and stages the artefacts
rem in deploy\ with the same layout as the CI Pi bundle.

echo Building frontend...
cd frontend
call npm ci
if errorlevel 1 exit /b 1
call npm run build
if errorlevel 1 exit /b 1
cd ..

echo Vetting backend...
cd backend
go vet ./...
if errorlevel 1 exit /b 1

echo Building backend for Raspberry Pi 2B (linux/arm v7)...
set GOOS=linux
set GOARCH=arm
set GOARM=7
go build -o ../deploy/brightlight .
if errorlevel 1 exit /b 1
set GOOS=
set GOARCH=
set GOARM=
cd ..

echo Copying frontend build to deploy...
robocopy frontend\build deploy\frontend\build /MIR /NFL /NDL /NJH /NJS /NC /NS /NP

echo Copying ui-config to deploy...
robocopy backend\ui-config deploy\backend\ui-config /MIR /NFL /NDL /NJH /NJS /NC /NS /NP

echo Copying systemd unit and installer to deploy...
copy /Y packaging\brightlight.service deploy\ >nul
copy /Y packaging\install.sh deploy\ >nul

echo Done. Deploy directory matches the CI bundle layout, ready for scp to Pi.
