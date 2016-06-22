go-importconst
==============

Import constants defined in C-header-files for Go ,
and makes `const.go` without `cgo`.

It requres `gcc.exe`.

    Usage: go-importconst.exe PACKAGENAME mark(constant)...
      -d ... do not remove temporary file
      -c ... clean output-files
      <header.h> \"header.h\" ... append headers
      NAME             -> const NAME=%d
      NAME:TYPE        -> const NAME=TYPE(%d)
      NAME:TYPE:FORMAT -> const NAME=TYPE(FORMAT)
    creates these files.
       -> ./makeconst.c (temporary)
       -> ./a.exe (temporary)
       -> ./const.go
    gcc and go-fmt are required.

Example-1
---------

### Commandline

    go-importconst ^
        conio ^
        CTRL_CLOSE_EVENT ^
        CTRL_LOGOFF_EVENT ^
        CTRL_SHUTDOWN_EVENT ^
        CTRL_C_EVENT ^
        ENABLE_ECHO_INPUT ^
        ENABLE_PROCESSED_INPUT ^
        STD_INPUT_HANDLE:uint32:0x%%X ^
        STD_OUTPUT_HANDLE:uint32:0x%%X

### Output (`const.go`)

    package conio

    const CTRL_CLOSE_EVENT = 2
    const CTRL_LOGOFF_EVENT = 5
    const CTRL_SHUTDOWN_EVENT = 6
    const CTRL_C_EVENT = 0
    const ENABLE_ECHO_INPUT = 4
    const ENABLE_PROCESSED_INPUT = 1
    const STD_INPUT_HANDLE = uint32(0xFFFFFFF6)
    const STD_OUTPUT_HANDLE = uint32(0xFFFFFFF5)

Example-2
---------

### Commandline

    go-importconst -d "stddef.h" "<stdlib.h>" main NULL

### Temporary-file (`makeconst.c`)

    #include <stdio.h>
    #include <windows.h>
    #include <stdlib.h>
    #include "stddef.h"

    #define d(n) printf("const " #n "=%d\n",n)
    #define s(n) printf("const " #n "=\"%s\"\n",n)
    #define u32x(n) printf("const " #n "=uint32(0x%08X)\n",n)
    #define up(n) printf("const " #n "=uintptr(%d)\n",n)
    #define MAKECONST_NULL(n) printf("const " #n "=%d\n",n)

    int main()
    {
        printf("package main\n\n");
        MAKECONST_NULL(NULL);
        return 0;
    }


### Output (`const.go`)

    package main

    const NULL = 0
