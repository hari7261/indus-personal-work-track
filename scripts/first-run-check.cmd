@echo off
setlocal
set SCRIPT_DIR=%~dp0
powershell -ExecutionPolicy Bypass -File "%SCRIPT_DIR%first-run-check.ps1"
if %ERRORLEVEL% neq 0 (
  echo.
  echo First-run check failed. Please review the messages above.
  exit /b %ERRORLEVEL%
)
echo.
echo First-run check passed.
endlocal
