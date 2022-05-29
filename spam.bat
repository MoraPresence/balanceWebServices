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
echo %Counter%
goto :WORK

:onTarget
echo %Counter%
goto :WORK
 
exit /b