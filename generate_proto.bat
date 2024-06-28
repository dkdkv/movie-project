@echo off
setlocal enabledelayedexpansion

:: Check if protoc is installed
where protoc >nul 2>nul
if %errorlevel% neq 0 (
    echo Error: protoc is not installed or not in PATH.
    echo Please install protoc and add it to your PATH.
    echo Visit https://github.com/protocolbuffers/protobuf/releases
    echo Download the appropriate zip file for your system (e.g., protoc-24.3-win64.zip^)
    echo Extract it and add the 'bin' directory to your PATH.
    exit /b 1
)

:: Set the root directory of your project
set PROJ_ROOT=%CD%

:: Set the protoc include paths
set PROTO_INCLUDE=-I"%PROJ_ROOT%\proto" -I"%GOPATH%\src"

:: Generate code for movie.proto
echo Generating code for: movie.proto

:: Generate Go code
protoc %PROTO_INCLUDE% --go_out="%OUT_DIR%" --go_opt=paths=source_relative --go-grpc_out="%OUT_DIR%" --go-grpc_opt=paths=source_relative "%PROJ_ROOT%\proto\movie\movie.proto"
if %errorlevel% neq 0 goto :error

:: Generate gRPC-Gateway code
protoc %PROTO_INCLUDE% --grpc-gateway_out="%OUT_DIR%" --grpc-gateway_opt=paths=source_relative "%PROJ_ROOT%\proto\movie\movie.proto"
if %errorlevel% neq 0 goto :error

:: Generate OpenAPI (Swagger) documentation
protoc %PROTO_INCLUDE% --openapiv2_out="%OUT_DIR%" --openapiv2_opt=allow_merge=true,merge_file_name=api "%PROJ_ROOT%\proto\movie\movie.proto"
if %errorlevel% neq 0 goto :error

echo Protocol buffer generation complete.
goto :eof

:error
echo Failed to generate protocol buffer code.
exit /b 1