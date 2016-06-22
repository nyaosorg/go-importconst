package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const CSOURCE = "makeconst.c"
const AEXE = "a.exe"
const CC = "gcc.exe"
const GOSOURCE = "const.go"

var debug = false

var packagename string
var headers []string
var stdheaders = []string{"<stdio.h>", "<windows.h>"}
var names = []string{}

var macros = map[string][2]string{}

func parse() bool {
	for _, arg1 := range os.Args[1:] {
		if arg1 == "-d" {
			debug = true
		} else if arg1 == "-c" {
			os.Remove(CSOURCE)
			os.Remove(GOSOURCE)
			return false
		} else if strings.HasSuffix(arg1, ".h") {
			headers = append(headers, arg1)
		} else if strings.HasSuffix(arg1, ".h>") {
			stdheaders = append(stdheaders, arg1)
		} else if p := strings.Split(arg1, ":"); len(p) == 3 {
			macros[p[0]] = [2]string{p[1], p[2]}
		} else if packagename == "" {
			packagename = arg1
		} else if strings.ContainsRune(arg1, '(') {
			names = append(names, arg1)
		} else {
			macros[arg1] = [2]string{"", "%d"}
		}
	}
	return true
}

func make_csource(csrcname string) {
	fd, err := os.Create(csrcname)
	if err != nil {
		fmt.Fprintf(fd, "%s: can not create makeconst.c\n", os.Args[0])
		return
	}
	defer fd.Close()

	for _, header1 := range stdheaders {
		fmt.Fprintf(fd, "#include %s\n", header1)
	}
	for _, header1 := range headers {
		fmt.Fprintf(fd, "#include \"%s\"\n", header1)
	}
	fmt.Fprintln(fd, ``)
	fmt.Fprintln(fd, `#define d(n) printf("const " #n "=%d\n",n)`)
	fmt.Fprintln(fd, `#define s(n) printf("const " #n "=\"%s\"\n",n)`)
	fmt.Fprintln(fd, `#define u32x(n) printf("const " #n "=uint32(0x%08X)\n",n)`)
	fmt.Fprintln(fd, `#define up(n) printf("const " #n "=uintptr(%d)\n",n)`)
	for key, val := range macros {
		fmt.Fprintf(fd, `#define MAKECONST_%s(n) printf("const " #n "=`, key)

		format := strings.Replace(val[1], `"`, `\"`, -1)
		if val[0] != "" {
			fmt.Fprintf(fd, `%s(%s)`, val[0], format)
		} else {
			fmt.Fprintf(fd, "%s", format)
		}
		fmt.Fprintln(fd, `\n",n)`)
	}
	fmt.Fprintln(fd, ``)
	fmt.Fprintln(fd, `int main()`)
	fmt.Fprintln(fd, `{`)
	fmt.Fprintln(fd, `    printf("package `+packagename+`\n\n");`)

	for _, name1 := range names {
		fmt.Fprintf(fd, "    %s;\n", name1)
	}
	for name1, _ := range macros {
		fmt.Fprintf(fd, "    MAKECONST_%s(%s);\n", name1, name1)
	}
	fmt.Fprintln(fd, "    return 0;\n}\n")
}

func compile() error {
	var gcc exec.Cmd
	gcc.Args = []string{
		CC,
		CSOURCE,
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
	const_c, err := os.Create(GOSOURCE)
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
		GOSOURCE,
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
	os.Remove(AEXE)
	err = gofmt()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if !parse() {
		return
	}
	if len(names) <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s PACKAGENAME mark(constant)...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "  -d ... do not remove temporary file")
		fmt.Fprintln(os.Stderr, "  -c ... clean output-files")
		fmt.Fprintln(os.Stderr, "  <header.h> \"header.h\" ... append headers")
		fmt.Fprintln(os.Stderr, "  d(NAME) ... const NAME=%d")
		fmt.Fprintln(os.Stderr, "  s(NAME) ... const NAME=\"%s\"")
		fmt.Fprintln(os.Stderr, "  u32x(NAME) ... const NAME=uint32(%08X)")
		fmt.Fprintln(os.Stderr, "  up(NAME) ... const NAME=uintptr(%d)")
		fmt.Fprintln(os.Stderr, "  NAME:TYPE:FORMAT ... const NAME=TYPE(FORMAT)")
		fmt.Fprintln(os.Stderr, "creates these files.")
		fmt.Fprintln(os.Stderr, "   -> ./makeconst.c (temporary)")
		fmt.Fprintln(os.Stderr, "   -> ./a.exe (temporary)")
		fmt.Fprintln(os.Stderr, "   -> ./const.go")
		fmt.Fprintln(os.Stderr, "gcc and go-fmt are required.")
		return
	}
	make_csource(CSOURCE)
	if !debug {
		defer os.Remove(CSOURCE)
	}

	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}
