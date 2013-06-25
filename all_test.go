package main

import (
    . "launchpad.net/gocheck"
    "testing"
	"path/filepath"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})


func (suite *MySuite) TestPathRestrict(c *C) {
	c.Check(filepath.Clean("foo"+"/"), Equals, "foo")
	c.Check(filepath.Clean("foo/"+"/"), Equals, "foo")
}
