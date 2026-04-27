@echo off
echo Building frontend...
cd frontend
call npm run build
cd ..

echo Building backend for Linux ARM (Raspberry Pi)...
cd backend
set GOOS=linux
set GOARCH=arm
set GOARM=5
go build -o ../deploy/brightlight
set GOOS=
set GOARCH=
set GOARM=
cd ..

echo Copying frontend build to deploy...
robocopy frontend\build deploy\frontend\build /MIR /NFL /NDL /NJH /NJS /NC /NS /NP

echo Done. Deploy directory is ready for scp to Pi.
