@setlocal
@set GOARCH=386
@set "PROMPT=MAKE>"
@call :"%1"
@endlocal
@exit /b

:""
    go fmt
    go build
    exit /b

:"test-dos"
    go-importconst -d -p dos ^
        FILE_ATTRIBUTE_NORMAL ^
        FILE_ATTRIBUTE_REPARSE_POINT ^
        FILE_ATTRIBUTE_HIDDEN ^
        CP_THREAD_ACP ^
        MOVEFILE_REPLACE_EXISTING ^
        MOVEFILE_COPY_ALLOWED ^
        MOVEFILE_WRITE_THROUGH
    exit /b
:"test-conio"
    go-importconst ^
        -d -p dos conio.h ^
        CTRL_CLOSE_EVENT ^
        CTRL_LOGOFF_EVENT ^
        CTRL_SHUTDOWN_EVENT ^
        CTRL_C_EVENT ^
        ENABLE_ECHO_INPUT ^
        ENABLE_PROCESSED_INPUT ^
        STD_INPUT_HANDLE ^
        STD_OUTPUT_HANDLE
    exit /b
:"test"
    go-importconst sample.h AHAHA IHIHI UFUFU
    exit /b
:"clean"
    go-importconst -c
    exit /b
