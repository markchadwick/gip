package main

import (
	"fmt"
	"strings"
)

type Scm struct {
	Name     string
	checkout func(url, branch, dst string) []string
}

var GitScm = &Scm{
	Name: "git",
	checkout: func(url, branch, dst string) []string {
		if branch == "" {
			branch = "master"
		}
		return []string{"clone", url, "--branch", branch, dst}
	},
}

func ScmFromPath(path string) (*Scm, error) {
	switch strings.Split(path, "/")[0] {
	default:
		return nil, fmt.Errorf("Could not find provider for %s", path)
	}
}

func (s *Scm) Checkout(url, branch, dst string) (string, []string) {
	return s.Name, s.checkout(url, branch, dst)
}

func (s *Scm) String() string {
	return s.Name
}

func pathParts(depth int) func(string) string {
	return func(path string) string {
		parts := strings.Split(path, "/")

		length := depth
		if length > len(parts) {
			length = len(parts)
		}
		return strings.Join(parts[:length], "/")
	}
}
