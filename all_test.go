package main

import (
	. "launchpad.net/gocheck"
	"path/filepath"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (suite *MySuite) TestPathRestrictBehaviourChecks(c *C) {
	c.Check(filepath.Clean("foo"+"/"), Equals, "foo")
	c.Check(filepath.Clean("foo/"+"/"), Equals, "foo")

	a, _ := filepath.Split("foo/a")
	c.Check(a, Equals, "foo/")
	c.Check(filepath.Clean(a), Equals, "foo")
}

func (suite *MySuite) TestCreateSpoolFilePath(c *C) {
	var testsym = []struct {
		spooldir string
		filename string
		ok       bool
		result   string
	}{
		{"/tmp", "foo", true, "/tmp/foo"},
		{"/tmp", "./foo", true, "/tmp/foo"},
		{"/tmp", "foo/a", false, ""},
		{"/tmp", "foo/b", false, ""},
		{"/tmp", "foo/../bar/1", false, ""},
		{"/tmp", "foo/../bar", true, "/tmp/bar"},
		{"/tmp", "/../foo", false, "/foo"},
		{"/tmp", "../foo", false, "/foo"},
		{"/tmp", "/../foo", false, "/foo"},
		{"/tmp", "/.././foo", false, "/foo"},
	}

	for _, sym := range testsym {
		test_result, test_ok := createSpoolFilePath(sym.spooldir, sym.filename)
		c.Check(test_ok, Equals, sym.ok, Commentf("Checked: %s + %s => %s", sym.spooldir, sym.filename, test_result))
		if sym.ok == true || (sym.ok == false && sym.result != "") {
			c.Check(test_result, Equals, sym.result)
		}
	}
}
