@echo off

cd %~dp0

set version=1.0.0b3

go build -o build\pvm.exe .
xcopy Scripts build\Scripts /I /Q /Y

cd build
tar.exe -a -c -f pvm-%version%.zip Scripts pvm.exe
