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
    creates these files.
       -> ./makeconst.c (temporary)
       -> ./a.exe (temporary)
       -> ./const.go
    gcc and go-fmt are required.

Example
-------

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

### Output

    package conio

    const CTRL_CLOSE_EVENT = 2
    const CTRL_LOGOFF_EVENT = 5
    const CTRL_SHUTDOWN_EVENT = 6
    const CTRL_C_EVENT = 0
    const ENABLE_ECHO_INPUT = 4
    const ENABLE_PROCESSED_INPUT = 1
    const STD_INPUT_HANDLE = uint32(0xFFFFFFF6)
    const STD_OUTPUT_HANDLE = uint32(0xFFFFFFF5)
