package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
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

	nLeafNodes, nInternalNodes := 0, 0
	dummyRoot.dfs(func(n *node) {
		// don't count the dummy root as a file or a dir
		if n == dummyRoot {
			return
		}

		if len(n.children) > 0 {
			nInternalNodes = nInternalNodes + 1
		} else {
			nLeafNodes = nLeafNodes + 1
		}
	})

	root := dummyRoot
	// if the dummy root only has a single child, we can use that
	// as the root for printing instead
	for len(root.children) == 1 {
		nextRoot := root.children[0]
		nextRoot.parent = nil
		nextRoot.name = path.Join(root.name, nextRoot.name)
		root = nextRoot
	}

	tree := root.PrintAsTree()
	itemCounts := fmt.Sprintf("%d %s, %d %s", nInternalNodes, directories(nInternalNodes), nLeafNodes, files(nLeafNodes))
	return fmt.Sprintf("%s\n%s\n", tree, itemCounts)
}

func directories(n int) string {
	if n == 1 {
		return "directory"
	}
	return "directories"
}

func files(n int) string {
	if n == 1 {
		return "file"
	}
	return "files"
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

func (n *node) dfs(visit func(*node)) {
	stack := []*node{n}

	pop := func() (next *node) {
		next, stack = stack[len(stack)-1], stack[:len(stack)-1]
		return next
	}

	seen := make(map[string]bool)

	for len(stack) > 0 {
		next := pop()

		// a nodes path uniquely identifies it
		p := next.printPath()
		if seen[p] {
			continue
		}

		visit(next)

		stack = append(stack, next.children...)

		seen[p] = true
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

func (n *node) printPath() string {
	parents := n.findParents()
	pathParts := make([]string, len(parents)+1)
	for _, parent := range parents {
		pathParts = append(pathParts, parent.name)
	}
	pathParts = append(pathParts, n.name)
	return path.Join(pathParts...)
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
