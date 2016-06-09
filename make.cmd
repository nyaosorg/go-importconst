@echo off
setlocal
set GOARCH=386

if not "%1" == "" goto %1
    go fmt importconst.go
    go build
    goto end
:test-null
    go-importconst -d "stddef.h" "<stdlib.h>" main d(NULL)
    goto end
:test-dos
    go-importconst dos ^
        d(FILE_ATTRIBUTE_NORMAL) ^
        d(FILE_ATTRIBUTE_REPARSE_POINT) ^
        d(FILE_ATTRIBUTE_HIDDEN) ^
        d(CP_THREAD_ACP) ^
        d(MOVEFILE_REPLACE_EXISTING) ^
        d(MOVEFILE_COPY_ALLOWED) ^
        d(MOVEFILE_WRITE_THROUGH)
    goto end
:test-conio
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
    goto end
:clean
    go-importconst -c
    goto end
:end

endlocal
