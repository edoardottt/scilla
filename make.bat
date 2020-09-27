@echo off

SET ARG=%1
SET TARGET=.\build
SET BUILDARGS=-ldflags="-s -w" -gcflags="all=-trimpath=%GOPATH%\src" -asmflags="all=-trimpath=%GOPATH%\src"

IF "%ARG%"=="" (
  go build -o .\scilla.exe
  GOTO Done
)

IF "%ARG%"=="windows" (
  CALL :Windows
  GOTO Done
)

IF "%ARG%"=="linux" (
  CALL :Linux
  GOTO Done
)

IF "%ARG%"=="update" (
  CALL :Update
  GOTO Done
)

IF "%ARG%"=="fmt" (
  CALL :Fmt
  GOTO Done
)

IF "%ARG%"=="remod" (
  del go.mod
  del go.sum
  go mod init github.com/edoardottt/scilla
  go get
  GOTO Done
)

IF "%ARG%"=="all" (
  CALL :Fmt
  CALL :Update
  CALL :Remod
  CALL :Test
  CALL :Linux
  CALL :Windows
  GOTO Done
)

IF "%ARG%"=="clean" (
  del /F /Q %TARGET%\*.*
  go clean ./...
  echo Done.
  GOTO Done
)

IF "%ARG%"=="test" (
  CALL :Test
  GOTO Done
)

GOTO Done

:Test
set GO111MODULE=on
set CGO_ENABLED=0
echo Testing ...
go test -v ./...
echo Done
EXIT /B 0

:Fmt
set GO111MODULE=on
echo Formatting ...
go fmt ./...
echo Done.
EXIT /B 0

:Update
set GO111MODULE=on
echo Updating ...
go get -u
go mod tidy -v
echo Done.
EXIT /B 0

:Linux
set GOOS=linux
set GOARCH=amd64
set GO111MODULE=on
set CGO_ENABLED=0
echo Building for %GOOS% %GOARCH% ...
set DIR=%TARGET%\scilla-%GOOS%-%GOARCH%
mkdir %DIR% 2> NUL
go build %BUILDARGS% -o %DIR%\scilla
set GOARCH=386
echo Building for %GOOS% %GOARCH% ...
set DIR=%TARGET%\scilla-%GOOS%-%GOARCH%
mkdir %DIR% 2> NUL
go build %BUILDARGS% -o %DIR%\scilla
echo Done.
EXIT /B 0

:Windows
set GOOS=windows
set GOARCH=amd64
set GO111MODULE=on
set CGO_ENABLED=0
echo Building for %GOOS% %GOARCH% ...
set DIR=%TARGET%\scilla-%GOOS%-%GOARCH%
mkdir %DIR% 2> NUL
go build %BUILDARGS% -o %DIR%\scilla.exe
set GOARCH=386
echo Building for %GOOS% %GOARCH% ...
set DIR=%TARGET%\scilla-%GOOS%-%GOARCH%
mkdir %DIR% 2> NUL
go build %BUILDARGS% -o %DIR%\scilla.exe
echo Done.
EXIT /B 0

:Done