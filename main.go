package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(printScannerAsTree(scanner))
}

func printScannerAsTree(s *bufio.Scanner) string {
	dummyRoot := &node{
		name:     ".",
		parent:   nil,
		children: []*node{},
	}
	for s.Scan() {
		path := s.Text()
		dummyRoot.insert(path)
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	root := dummyRoot
	// if the dummy root only has a single child, we can use that
	// as the root for printing instead
	if len(root.children) == 1 {
		root = root.children[0]
		root.parent = nil
	}

	return root.PrintAsTree()
}

type node struct {
	parent   *node
	children []*node
	name     string
}

const (
	verticalPipe             = '│'
	horizontalPipe           = '─'
	cornerPipe               = '└'
	verticalPipeWithOffshoot = '├'
)

// PrintAsTree prints the node & all it's children as a `tree`-style tree
func (n *node) PrintAsTree() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n", n.name))
	n.printAsTreeHelper(&sb)
	return sb.String()
}

func (n *node) printAsTreeHelper(sb *strings.Builder) {
	for _, child := range n.children {
		for _, parent := range child.findParents() {
			if parent.isRoot() {
				continue
			}

			var connChar rune
			if parent.isLastChild() {
				connChar = ' '
			} else {
				connChar = verticalPipe
			}
			sb.WriteString(fmt.Sprintf("%c%s", connChar, spaces(3)))
		}

		var connChar rune
		if child.isLastChild() {
			connChar = cornerPipe
		} else {
			connChar = verticalPipeWithOffshoot
		}

		sb.WriteString(fmt.Sprintf("%c%c%c %s\n", connChar, horizontalPipe, horizontalPipe, child.name))
		child.printAsTreeHelper(sb)
	}
}

// insert inserts all nodes represented by the supplied path. A node is added for each segment.
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

// findParents returns an array containing all the nodes parents. The parent nodes are
// returned in order of highest to lowest, and the root node is skipped. That is, the
// first node will be the parent node closest to the root.
func (n *node) findParents() []*node {
	parents := []*node{}
	current := n
	for current.parent != nil {
		parents = append([]*node{current.parent}, parents...)
		current = current.parent
	}
	return parents
}

// findChildWithName finds a child with a name matching the supplied name inside the node's
// array of child nodes
func (n *node) findChildWithName(name string) *node {
	for _, child := range n.children {
		if child.name == name {
			return child
		}
	}
	return nil
}

// isRoot returns true if the node has no parent
func (n *node) isRoot() bool {
	return n.parent == nil
}

// isLastChild returns true if n is the final node in the parent node's array of children
func (n *node) isLastChild() bool {
	if n.parent == nil {
		return false
	}
	return n.position() == len(n.parent.children)-1
}

// position returns the index of the current node in the parent node's array of children
func (n *node) position() int {
	i := n.parent.indexOf(n)
	if i == -1 {
		panic("n is not a child of its parent")
	}
	return i
}

// indexOf returns the index of `target` in the `children` array of `n`, if itexists. Otherwise,
// it returns -1
func (n *node) indexOf(target *node) int {
	for i, child := range n.children {
		if child == target {
			return i
		}
	}
	return -1
}

// spaces returns a string of length n containing only space characters
func spaces(n int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = ' '
	}
	return string(s)
}

func newNode(name string, parent *node, children []*node) *node {
	return &node{
		parent,
		children,
		name,
	}
}
