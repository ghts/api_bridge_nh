@echo off

SET CGO_ENABLED=1
SET GOARCH=386
SET DEPENDENCY_PATH=%GOPATH%\src\github.com\ghts\ghts_dependency
SET GCC_PATH=%DEPENDENCY_PATH%\ruby_devkit_32
SET C_INCLUDE_PATH=%DEPENDENCY_PATH%\build_dep_32\include
SET LIBRARY_PATH=%DEPENDENCY_PATH%\build_dep_32\lib
SET DLL_PATH=%DEPENDENCY_PATH%\build_dep_32\bin
SET NH_OpenAPI_PATH=%DEPENDENCY_PATH%\NH_OpenAPI

SET OLDPATH=%PATH%
SET PATH=%GCC_PATH%\bin;%GCC_PATH%\mingw\bin;%DLL_PATH%;%NH_OpenAPI_PATH%;%PATH%

REM Bootstrapping cross-compile
REM go tool dist install -v runtime
REM go install -v -a std

cls
go build api_bridge_nh.go

SET PATH=%OLDPATH%