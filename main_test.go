package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func assertNotNil(t *testing.T, n *node) {
	if n == nil {
		t.Error("exptected to find a node, instead got nil")
	}
}

func assertIs(t *testing.T, n1, n2 *node) {
	assertNotNil(t, n1)
	if n1 != n2 {
		t.Errorf("exptected to find a node with name \"%s\", instead found node with name \"%s\"", n2.name, n1.name)
	}
}

func TestFindChildWithName(t *testing.T) {
	root := newNode("root", nil, []*node{})
	n1 := newNode("node1", root, []*node{})
	n2 := newNode("node2", root, []*node{})

	root.children = append(root.children, n1, n2)

	n := root.findChildWithName("node1")
	assertIs(t, n, n1)

	n = root.findChildWithName("node2")
	assertIs(t, n, n2)
}

func TestInsertDoesNotDuplicateChildren(t *testing.T) {
	root := newNode("root", nil, []*node{})
	root.insert("src")
	root.insert("src")
	if len(root.children) != 1 {
		t.Errorf("expected root to have 1 child, got %d children", len(root.children))
	}

	root = newNode("root", nil, []*node{})
	root.insert("src/dir1")
	root.insert("src/dir1")

	srcNode := root.findChildWithName("src")
	assertNotNil(t, srcNode)
	if len(srcNode.children) != 1 {
		t.Errorf("expected srcNode to have 1 child, got %d children", len(root.children))
	}
}

func TestInsertAddsChildrenWithCommonParent(t *testing.T) {
	root := newNode("root", nil, []*node{})
	root.insert("src/foo")
	root.insert("src/bar")

	srcNode := root.findChildWithName("src")
	assertNotNil(t, srcNode)
	assertNotNil(t, srcNode.findChildWithName("foo"))
	assertNotNil(t, srcNode.findChildWithName("bar"))
}

func TestFixtures(t *testing.T) {
	files, err := ioutil.ReadDir("fixtures")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			t.Error("expected everything in fixtures/ to be a directory")
		}

		inFile, err := os.Open(path.Join("fixtures", f.Name(), "input"))
		if err != nil {
			t.Fatal(err)
		}

		wantBytes, err := ioutil.ReadFile(path.Join("fixtures", f.Name(), "want"))
		if err != nil {
			t.Fatal(err)
		}

		want := strings.TrimSpace(string(wantBytes))
		got := strings.TrimSpace(printScannerAsTree(bufio.NewScanner(inFile)))
		if want != got {
			t.Errorf("input for fixture %s did not match output.\nwanted:\n%s\ngot:\n%s", f.Name(), want, got)
		}
	}
}
