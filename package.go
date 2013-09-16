package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Package struct {
	Path   string
	scm    *Scm
	url    string
	branch string
}

func PackageFromImport(i string) (*Package, error) {
	provider := strings.Split(i, "/")[0]
	return PackageFromProvider(provider, i)
}

func PackageFromProvider(provider string, i string) (*Package, error) {
	pkg := &Package{
		Path: i,
	}
	// TODO in SCM now
	switch provider {
	default:
		return nil, fmt.Errorf("Can't find provider '%s'", provider)
	case "github.com":
		pkg.scm = GitScm
		pkg.url = fmt.Sprintf("https://%s", i)
	}
	return pkg, nil
}

func (p *Package) Nme() string {
	return p.Path[strings.LastIndex(p.Path, "/")+1:]
}

func (p *Package) Checkout(dst string) error {
	cmd, args := p.scm.Checkout(p.url, p.branch, dst)
	return inDir(dst, func() error {
		return p.exec(cmd, args)
	})
}

func (p *Package) exec(name string, args []string) error {
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)

	cmd := exec.Command(name, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Printf("exec: %s %s", name, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error running %s %s\n%s\n\nstdout:\n%s\nstderr:\n%s\n",
			name, strings.Join(args, " "), err.Error(), stdout.String(),
			stderr.String())
	}
	return nil
}

func inDir(dir string, f func() error) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	if err := os.Chdir(dir); err != nil {
		return err
	}
	return f()
}
