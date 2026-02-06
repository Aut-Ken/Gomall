@echo off
REM Database Backup Script for Gomall (Windows)
REM 使用方法: backup.bat [backup_dir]

setlocal

set BACKUP_DIR=%~1
if "%BACKUP_DIR%"=="" set BACKUP_DIR=.\backups

set DATE=%date:~0,4%%date:~5,2%%date:~8,2%_%time:~0,2%%time:~3,2%
set DATE=%DATE: =0%
set BACKUP_FILE=%BACKUP_DIR%\gomall_%DATE%.sql

set DB_HOST=%GOMALL_DB_HOST:-localhost=localhost%
set DB_PORT=%GOMALL_DB_PORT:-3306=3306%
set DB_USER=%GOMALL_DB_USER:-root=root%
set DB_PASSWORD=%GOMALL_DB_PASSWORD:-password=password%
set DB_NAME=%GOMALL_DB_NAME:-Gomall=Gomall%

mkdir "%BACKUP_DIR%" 2>nul

echo 开始备份数据库: %DB_NAME%
echo 备份文件: %BACKUP_FILE%

mysqldump -h "%DB_HOST%" -P "%DB_PORT%" -u "%DB_USER%" -p"%DB_PASSWORD%" --single-transaction --routines --triggers "%DB_NAME%" > "%BACKUP_FILE%"

if exist "%BACKUP_FILE%" (
    for %%I in ("%BACKUP_FILE%") do echo 备份成功! 文件大小: %%~zI bytes
) else (
    echo 备份失败!
    exit /b 1
)

echo.
echo 备份完成!
endlocal
