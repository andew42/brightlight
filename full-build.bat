cd ui2
start npm run build
cd ..
set GOOS=linux
set GOARCH=arm
set GOARM=5
go build
set GOOS=
set GOARCH=
set GOARM=
