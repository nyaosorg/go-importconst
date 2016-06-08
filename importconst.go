package main

import (
	"fmt"
	"os"
	"os/exec"
)

const CSOURCE = "makeconst.c"
const AEXE = "a.exe"
const CC = "gcc.exe"

func make_csource(csrcname string) {
	fd, err := os.Create(csrcname)
	if err != nil {
		fmt.Fprintf(fd, "%s: can not create makeconst.c\n", os.Args[0])
		return
	}
	defer fd.Close()

	fmt.Fprintln(fd, `#include <stdio.h>`)
	fmt.Fprintln(fd, `#include <windows.h>`)
	fmt.Fprintln(fd, ``)
	fmt.Fprintln(fd, `#define d(n) printf("const " #n "=%d\n",n)`)
	fmt.Fprintln(fd, `#define s(n) printf("const " #n "=\"%s\"\n",n)`)
	fmt.Fprintln(fd, `#define u32x(n) printf("const " #n "=uint32(0x%08X)\n",n)`)
	fmt.Fprintln(fd, ``)
	fmt.Fprintln(fd, `int main()`)
	fmt.Fprintln(fd, `{`)
	fmt.Fprintln(fd, `    printf("package `+os.Args[1]+`\n\n");`)

	for _, arg1 := range os.Args[2:] {
		fmt.Fprintf(fd, "    %s;\n", arg1)
	}
	fmt.Fprintln(fd, "    return 0;\n}\n")
}

func compile() error {
	var gcc exec.Cmd
	gcc.Args = []string{
		CC,
		"makeconst.c",
	}
	gcc.Path = gcc.Args[0]
	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr
	return gcc.Run()
}

func aexe() error {
	var aexe exec.Cmd
	aexe.Args = []string{
		AEXE,
	}
	aexe.Path = aexe.Args[0]
	const_c, err := os.Create("const.go")
	if err != nil {
		return err
	}
	defer const_c.Close()
	aexe.Stdout = const_c
	aexe.Stderr = os.Stderr
	return aexe.Run()
}

func gofmt() error {
	var gofmt exec.Cmd
	gofmt.Args = []string{
		"go",
		"fmt",
		"const.go",
	}
	gofmt.Path = gofmt.Args[0]
	gofmt.Stdout = os.Stdout
	gofmt.Stderr = os.Stderr
	return gofmt.Run()
}

func main1() error {
	err := compile()
	if err != nil {
		return err
	}
	err = aexe()
	if err != nil {
		return err
	}
	defer os.Remove(AEXE)
	err = gofmt()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: %s PACKAGENAME mark(constant)...")
		fmt.Fprintln(os.Stderr, "   d(NAME) ... const NAME=%d")
		fmt.Fprintln(os.Stderr, "   s(NAME) ... const NAME=\"%s\"")
		fmt.Fprintln(os.Stderr, "   u32x(NAME) ... const NAME=uint32(%08X)")
		fmt.Fprintln(os.Stderr, "creates these files.")
		fmt.Fprintln(os.Stderr, "   -> ./makeconst.c (temporary)")
		fmt.Fprintln(os.Stderr, "   -> ./a.exe (temporary)")
		fmt.Fprintln(os.Stderr, "   -> ./const.go")
		fmt.Fprintln(os.Stderr, "gcc and go-fmt are required.")
		return
	}

	make_csource(CSOURCE)
	defer os.Remove(CSOURCE)

	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}
