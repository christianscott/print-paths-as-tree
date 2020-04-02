package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type node struct {
	parent   *node
	children []*node
	name     string
}

func newNode(name string, parent *node, children []*node) *node {
	return &node{
		parent,
		children,
		name,
	}
}

func (n *node) insert(path string) {
	current := n
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		next := current.findChildWithName(segment)
		// if there is no node at this level with a name matching the current
		// segment, create a new node and add it as a child of "current"
		if next == nil {
			next = &node{
				name:     segment,
				parent:   current,
				children: []*node{},
			}
			current.children = append(current.children, next)
		}

		current = next
	}
}

func (n *node) findChildWithName(name string) *node {
	for _, child := range n.children {
		if child.name == name {
			return child
		}
	}
	return nil
}

func (n *node) findParents() []*node {
	parents := []*node{}
	current := n
	for current.parent != nil && !current.parent.isRoot() {
		parents = append([]*node{current.parent}, parents...)
		current = current.parent
	}
	return parents
}

func (n *node) isRoot() bool {
	return n.parent == nil
}

func (n *node) indexOf(target *node) int {
	for i, child := range n.children {
		if child == target {
			return i
		}
	}
	return -1
}

func (n *node) position() int {
	i := n.parent.indexOf(n)
	if i == -1 {
		panic("n is not a child of its parent")
	}
	return i
}

func (n *node) isLastChild() bool {
	if n.parent == nil {
		return false
	}
	return n.position() == len(n.parent.children)-1
}

func (n *node) printAsTree() string {
	var sb strings.Builder
	sb.WriteString(".\n")
	printAsTreeHelper(&sb, n)
	return sb.String()
}

const (
	verticalPipe             = '│'
	horizontalPipe           = '─'
	cornerPipe               = '└'
	verticalPipeWithOffshoot = '├'
)

func printAsTreeHelper(sb *strings.Builder, n *node) {
	for _, child := range n.children {
		for _, parent := range child.findParents() {
			var connChar rune
			if parent.isLastChild() {
				connChar = ' '
			} else {
				connChar = verticalPipe
			}
			sb.WriteString(fmt.Sprintf("%c%s", connChar, spaces(len(parent.name)-1)))
		}

		var connChar rune
		if child.isLastChild() {
			connChar = cornerPipe
		} else {
			connChar = verticalPipeWithOffshoot
		}

		sb.WriteString(fmt.Sprintf("%c%c %s\n", connChar, horizontalPipe, child.name))
		printAsTreeHelper(sb, child)
	}
}

func spaces(n int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = ' '
	}
	return string(s)
}

func printScannerAsTree(s *bufio.Scanner) string {
	root := &node{
		name:     "root",
		parent:   nil,
		children: []*node{},
	}
	for s.Scan() {
		path := s.Text()
		root.insert(path)
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	return root.printAsTree()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(printScannerAsTree(scanner))
}
