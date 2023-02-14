package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var (
	flagCSrc  = flag.String("csrc", "zconst.cpp", "c-source filename used temporally")
	flagCc    = flag.String("cc", "gcc", "c compiler command")
	flagGoSrc = flag.String("o", "zconst.go", "go-source-filename to output constants")
	flagClean = flag.Bool("c", false, "clean output")
	flagDebug = flag.Bool("d", false, "debug flag")
	flagNofmt = flag.Bool("nofmt", false, "do not execute go fmt (for debug)")
)

var packageName = os.Getenv("GOPACKAGE")

func makeCSource(csrcname string, headers []string, vars []string) {
	fd, err := os.Create(csrcname)
	if err != nil {
		fmt.Fprintf(fd, "%s: can not create %s\n", os.Args[0], csrcname)
		return
	}
	defer fd.Close()

	for _, header1 := range headers {
		fmt.Fprintf(fd, "#include %s\n", header1)
	}
	fmt.Fprint(fd, `
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
`)
	fmt.Fprintln(fd, `    printf("package `+packageName+`\n\n");`)
	fmt.Fprintln(fd, `    printf("// Code generated by go-importconst DO NOT EDIT.\n");`)

	for _, name1 := range vars {
		fmt.Fprintf(fd, "    p(\"%s\",%s);\n", name1, name1)
	}
	fmt.Fprintln(fd, "    return 0;\n}\n")
}

func compile() error {
	var cc exec.Cmd
	cc.Args = []string{
		*flagCc,
		*flagCSrc,
	}
	fn, err := exec.LookPath(*flagCc)
	if err != nil {
		return err
	}
	cc.Path = fn
	cc.Stdout = os.Stdout
	cc.Stderr = os.Stderr
	return cc.Run()
}

func nameOfExecutable() string {
	if runtime.GOOS == "windows" {
		return "a.exe"
	} else {
		return "a.out"
	}
}

func aexe() (string, error) {
	constC, err := os.Create(*flagGoSrc)
	if err != nil {
		return "", err
	}
	defer constC.Close()

	fname := nameOfExecutable()
	aexe := exec.Cmd{
		Args:   []string{fname},
		Path:   fname,
		Stdout: constC,
		Stderr: os.Stderr,
	}
	return fname, aexe.Run()
}

func gofmt() error {
	if *flagNofmt {
		return nil
	}
	var gofmt exec.Cmd
	gofmt.Args = []string{
		"go",
		"fmt",
		*flagGoSrc,
	}
	fn, err := exec.LookPath("go")
	if err != nil {
		return err
	}
	gofmt.Path = fn
	gofmt.Stdout = os.Stdout
	gofmt.Stderr = os.Stderr
	return gofmt.Run()
}

func readGoGenerateParameter() ([]string, error) {
	gofile := os.Getenv("GOFILE")
	if gofile == "" {
		return nil, errors.New("$GOFILE is not defined. Use `go generate`")
	}
	fd, err := os.Open(gofile)
	if err != nil {
		return nil, err
	}
	goline := os.Getenv("GOLINE")
	if goline == "" {
		return nil, errors.New("$GOLINE is not defined. Use `go generate`")
	}
	lnum, err := strconv.Atoi(goline)
	if err != nil {
		return nil, fmt.Errorf("$GOLINE: %s", err.Error())
	}
	sc := bufio.NewScanner(fd)
	var tokens []string
	for sc.Scan() {
		lnum--
		if lnum >= 0 {
			continue
		}
		text := sc.Text()
		if !strings.HasPrefix(text, "//") {
			break
		}
		fields := strings.Fields(text[2:])
		if len(fields) <= 0 {
			break
		}
		for _, arg1 := range fields {
			tokens = append(tokens, arg1)
		}
	}
	return tokens, nil
}

func remove(fn string) {
	fmt.Fprintln(os.Stderr, "rm", fn)
	os.Remove(fn)
}

func main1() error {
	if *flagClean {
		remove(*flagCSrc)
		remove(nameOfExecutable())
		remove(*flagGoSrc)
		return nil
	}
	goParams, err := readGoGenerateParameter()
	if err != nil {
		return err
	}
	headers := []string{"<cstdio>"}
	vars := make([]string, 0)
	for _, s := range goParams {
		if len(s) > 0 && s[0] == '<' {
			headers = append(headers, s)
		} else if strings.HasSuffix(s, ".h") {
			headers = append(headers, fmt.Sprintf(`"%s"`, s))
		} else {
			vars = append(vars, s)
		}
	}

	makeCSource(*flagCSrc, headers, vars)

	if !*flagDebug {
		defer remove(*flagCSrc)
	}

	err = compile()
	if err != nil {
		return err
	}
	fname, err := aexe()
	if err != nil {
		return err
	}
	remove(fname)
	err = gofmt()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
