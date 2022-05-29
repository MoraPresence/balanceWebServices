set /A Counter=0
:WORK
curl http://127.0.0.1:8080
timeout 1 > nul
set /A Counter+=1
if %Counter% equ 100 (
	goto :offTarget
)
if %Counter% equ 200 (
	goto :onTarget
)
goto :WORK
 
:offTarget
docker-compose.exe -f .\docker-compose.yml stop target-service01
goto :WORK

:onTarget
docker-compose.exe -f .\docker-compose.yml up -d target-service01
goto :WORK
 
exit /b