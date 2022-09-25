@echo off
setlocal enabledelayedexpansion
call build.bat
if not "%ERRORLEVEL%"=="0" ( goto:eof )
if "%1" == "" (
    call run.bat change
) else (
    call run.bat %*
)
if not "%ERRORLEVEL%"=="0" ( goto:eof )
