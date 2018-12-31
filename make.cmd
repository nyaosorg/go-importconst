@setlocal
@set GOARCH=386
@set "PROMPT=MAKE>"
@call :"%1"
@endlocal
@exit /b

:""
    go fmt importconst.go
    go build
    exit /b

:"test-null"
    go-importconst -d "stddef.h" "<stdlib.h>" main d(NULL)
    exit /b

:"test-dos"
    go-importconst dos ^
        d(FILE_ATTRIBUTE_NORMAL) ^
        d(FILE_ATTRIBUTE_REPARSE_POINT) ^
        d(FILE_ATTRIBUTE_HIDDEN) ^
        d(CP_THREAD_ACP) ^
        d(MOVEFILE_REPLACE_EXISTING) ^
        d(MOVEFILE_COPY_ALLOWED) ^
        d(MOVEFILE_WRITE_THROUGH)
    exit /b
:"test-conio"
    go-importconst ^
        conio ^
        d(CTRL_CLOSE_EVENT) ^
        d(CTRL_LOGOFF_EVENT) ^
        d(CTRL_SHUTDOWN_EVENT) ^
        d(CTRL_C_EVENT) ^
        d(ENABLE_ECHO_INPUT) ^
        d(ENABLE_PROCESSED_INPUT) ^
        u32x(STD_INPUT_HANDLE) ^
        u32x(STD_OUTPUT_HANDLE)
    exit /b
:"clean"
    go-importconst -c
    exit /b
