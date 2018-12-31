go-importconst
==============

Import constants defined in C-header-files for Go ,
and makes `const.go` without `cgo`.

It requires `gcc`.

    Usage:
        go-importconst {OPTIONS} {HEADERFILE(S)...} {CONSTANTS}...
            OR
        go run importconst.go {OPTIONS} {HEADERFILE(S)...} {CONSTANTS}...

    options:
      -d ... For debug, do not remove temporary file
      -p PACKAGE ... specify packagename(default: main)
    creates these files.
       -> ./makeconst.c (temporary)
       -> ./a.exe (temporary)
       -> ./const.go
    gcc and go are required.

Example
---------

sample.h

    #ifndef HOGE_H
    #  define AHAHA "ahaha"
    #  define IHIHI 12345
    #  define UFUFU 3.14159
    #endif

Commandline:

    $ go-importconst sample.h AHAHA IHIHI UFUFU

Output(`const.go`):
    package main

    const AHAHA = "ahaha"
    const IHIHI = 12345
    const UFUFU = 3.141590

A temporary file(`makeconst.c`):

    #include <cstdio>
    #include <windows.h>
    #include "sample.h"

    void p(const char *name,const char *s){
        printf("const %s=\"%s\"\n",name,s);
    }
    void p(const char *name,int n){
        printf("const %s=%d\n",name,n);
    }
    void p(const char *name,long n){
        printf("const %s=%ld\n",name,n);
    }
    void p(const char *name,unsigned long n){
        printf("const %s=%ld\n",name,n);
    }
    void p(const char *name,double n){
        printf("const %s=%lf\n",name,n);
    }

    int main()
    {
        printf("package main\n\n");
        p("AHAHA",AHAHA);
        p("IHIHI",IHIHI);
        p("UFUFU",UFUFU);
        return 0;
    }
