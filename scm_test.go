package main

import (
	"./github.com/markchadwick/spec"
	"testing"
)

func Test(t *testing.T) {
	spec.Run(t)
}

var _ = spec.Suite("path parts", func(c *spec.C) {
	f := pathParts(3)

	c.It("should return the properly lengthed path", func(c *spec.C) {
		c.Assert(f("github.com/user/project")).Equals("github.com/user/project")
		c.Assert(f("github.com/user/project/foo")).Equals("github.com/user/project")
	})

	c.It("should accept short paths", func(c *spec.C) {
		c.Assert(f("github.com/user")).Equals("github.com/user")
		c.Assert(f("github.com")).Equals("github.com")
		c.Assert(f("")).Equals("")
	})

	c.It("should ignore trailing slashes", func(c *spec.C) {
		c.Assert(f("github.com/user/project/")).Equals("github.com/user/project")
	})
})
