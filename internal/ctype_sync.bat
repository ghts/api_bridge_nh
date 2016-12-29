@echo off

call %GOPATH%\ghts_dependency\batch_scripts\32.bat

cls

SET OLDPATH=%PATH%
SET PATH=%GCC_PATH%\bin;%GCC_PATH%\mingw\bin;%DLL_PATH%;%NH_OpenAPI_PATH%;%PATH%

cd %PROJECT_ROOT%\ghts_connector_nh\internal

copy ctype.orig ctype_1.go
go tool cgo -godefs ctype_1.go > ctype_2.go
sed -e 's/int8/byte/g' ctype_2.go > ctype.go

del ctype_?.go

go fmt

SET PATH=%OLDPATH%