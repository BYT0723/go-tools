package main

import (
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	var (
		fs   = flag.NewFlagSet("global", flag.ExitOnError)
		args = os.Args[1:]
	)

	fs.Parse(args)

	subcmd := fs.Arg(0)

	switch subcmd {
	case "makegen":
		makegen(args)
	}
}

//go:embed makefile.tmpl
var makefiletmpl string

// makefile generate
func makegen(args []string) {
	var name, buildDir, importPkg, output string

	fs := flag.NewFlagSet("makegen", flag.ExitOnError)

	wd, _ := os.Getwd()

	fs.StringVar(&name, "name", filepath.Base(wd), "application name")
	fs.StringVar(&buildDir, "build-dir", "./bin", "build dir")
	fs.StringVar(&importPkg, "import-pkg", "main", "import value package")
	fs.StringVar(&output, "output", "Makefile", "out file")

	fs.Parse(args)

	t, err := template.New("makefiletmpl").Parse(makefiletmpl)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := t.Execute(f, map[string]any{
		"name":      name,
		"buildDir":  buildDir,
		"importPkg": importPkg,
	}); err != nil {
		panic(err)
	}
}
