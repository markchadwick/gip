package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Path struct {
	Root      string
	installed map[string]int
}

func OpenPath() (*Path, error) {
	p := &Path{
		installed: make(map[string]int),
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	p.Root = path.Join(cwd, ".gip")
	return p, nil
}

func (p *Path) init() error {
	return os.MkdirAll(p.Root, 0755)
}

func (p *Path) Install(pkg *Package) error {
	if _, ok := p.installed[pkg.Path]; ok {
		log.Printf("%s already installed", pkg.Path)
		return nil
	}
	p.installed[pkg.Path] = 0

	if err := p.download(pkg); err != nil {
		return err
	}

	return p.install(pkg)
}

func (p *Path) install(pkg *Package) (err error) {
	for _, imprt := range p.imports(p.srcPath(pkg.Path)) {
		log.Printf("%s ⇒ %s", pkg.Path, imprt)
		dep, err := PackageFromImport(imprt)
		if err != nil {
			return err
		}
		if err := p.Install(dep); err != nil {
			return err
		}
	}

	return p.build(pkg)
}

func (p *Path) build(pkg *Package) error {
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	path := p.srcPath(pkg.Path)

	cmd := exec.Command("go", "install")
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	env := []string{
		"PWD=" + path,
		"GOPATH=" + p.Root,
	}
	cmd.Env = append(os.Environ(), env...)
	log.Printf("exec: go install (in %s)", path)
	if err := inDir(path, func() error { return cmd.Run() }); true {
		return fmt.Errorf("Error running go install %s\n%s\n\nstdout:\n%s\nstderr:\n%s\n",
			path,
			err, stdout.String(),
			stderr.String())
	}

	return nil
}

func (p *Path) download(pkg *Package) (err error) {
	stat, err := os.Stat(p.srcPath(pkg.Path))
	if err == nil && stat.IsDir() {
		err = p.update(pkg)
	} else {
		err = p.checkout(pkg)
	}
	return nil
}

func (p *Path) update(pkg *Package) error {
	log.Printf("%s updating (TODO)", pkg.Path)
	return nil
}

func (p *Path) checkout(pkg *Package) error {
	log.Printf("%s checking out", pkg.Path)
	root := p.srcPath(pkg.Path)
	if err := os.MkdirAll(root, 0755); err != nil {
		return err
	}
	return pkg.Checkout(root)
}

func (p *Path) srcPath(pth string) string {
	return path.Join(p.Root, "src", pth)
}

func (p *Path) imports(root string) []string {
	all := make(chan string)
	go func() {
		defer close(all)
		filepath.Walk(root, importWalker(all))
	}()

	seen := make(map[string]int)
	imports := make([]string, 0)

	for i := range all {
		if _, ok := seen[i]; ok {
			continue
		}
		seen[i] = 0
		if isRemote(i) {
			imports = append(imports, i)
		}
	}

	return imports
}

func importWalker(imports chan string) filepath.WalkFunc {
	fset := token.NewFileSet()

	return func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") || info.IsDir() {
			return err
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, i := range f.Imports {
			quoted := i.Path.Value
			imports <- quoted[1 : len(quoted)-1]
		}

		return nil
	}
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func isRemote(path string) bool {
	slash := strings.Index(path, "/")
	dot := strings.Index(path, ".")

	return slash >= 0 && dot >= 0 && dot < slash
}
