go-importconst
==============

Import constants defined in C-header-files for Go ,
and makes `const.go` without `cgo`.

It requres `gcc.exe`.

    Usage: go-importconst.exe PACKAGENAME mark(constant)...
      -d ... do not remove temporary file
      -c ... clean output-files
      <header.h> "header.h" ... append headers
      d(NAME) ... const NAME=%d
      s(NAME) ... const NAME="%s"
      u32x(NAME) ... const NAME=uint32(%08X)
      up(NAME) ... const NAME=uintptr(%d)
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
        d(CTRL_CLOSE_EVENT) ^
        d(CTRL_LOGOFF_EVENT) ^
        d(CTRL_SHUTDOWN_EVENT) ^
        d(CTRL_C_EVENT) ^
        d(ENABLE_ECHO_INPUT) ^
        d(ENABLE_PROCESSED_INPUT) ^
        u32x(STD_INPUT_HANDLE) ^
        u32x(STD_OUTPUT_HANDLE)

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

    go-importconst -d "stddef.h" "<stdlib.h>" main d(NULL)

### Temporary-file (`makeconst.c`)

    #include <stdio.h>
    #include <windows.h>
    #include <stdlib.h>
    #include "stddef.h"

    #define d(n) printf("const " #n "=%d\n",n)
    #define s(n) printf("const " #n "=\"%s\"\n",n)
    #define u32x(n) printf("const " #n "=uint32(0x%08X)\n",n)

    int main()
    {
        printf("package main\n\n");
        d(NULL);
        return 0;
    }

### Output (`const.go`)

    package main

    const NULL = 0
