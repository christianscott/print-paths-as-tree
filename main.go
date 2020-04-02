package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

type node struct {
	id       int
	parent   *node
	children []*node
	name     string
}

var nextNodeID int

func newNode(name string, parent *node, children []*node) *node {
	id := nextNodeID
	nextNodeID = nextNodeID + 1
	return &node{
		id,
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
			next = newNode(
				segment,
				current,
				[]*node{},
			)
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

func (n *node) dfs(callback func(*node)) {
	seen := make(map[int]bool)

	stack := make([]*node, len(n.children))
	copy(stack, n.children)

	for len(stack) > 0 {
		var child *node
		child, stack = stack[len(stack)-1], stack[:len(stack)-1]

		if seen[child.id] {
			continue
		}

		callback(child)

		if len(child.children) > 0 {
			stack = append(stack, child.children...)
		}

		seen[child.id] = true
	}
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

func (n *node) getPath() string {
	parents := n.findParents()
	paths := make([]string, len(parents)+1)
	for _, parent := range parents {
		paths = append(paths, parent.name)
	}
	paths = append(paths, n.name)
	return path.Join(paths...)
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

func (n *node) hasChildren() bool {
	return len(n.children) > 0
}

func (n *node) printAsTree() string {
	var sb strings.Builder
	sb.WriteString(".\n")
	// for i, child := range n.children {
	// 	if i == len(n.children)-1 && len(child.children) == 0 {
	// 		sb.WriteString(fmt.Sprintf("%c%c %s\n", cornerPipe, horizontalPipe, child.name))
	// 	} else {
	// 		sb.WriteString(fmt.Sprintf("%c%c %s\n", verticalPipeWithOffshoot, horizontalPipe, child.name))
	// 	}
	// 	for j, childsChild := range child.children {
	// 		if j == len(child.children)-1 && len(childsChild.children) == 0 {
	// 			sb.WriteString(fmt.Sprintf("%c%s%c%c %s\n", verticalPipe, spaces(len(child.name)-1), cornerPipe, horizontalPipe, childsChild.name))
	// 		} else {
	// 			sb.WriteString(fmt.Sprintf("%c%s%c%c %s\n", verticalPipe, spaces(len(child.name)-1), verticalPipeWithOffshoot, horizontalPipe, childsChild.name))
	// 		}
	// 	}
	// }
	printAsTreeHelper(&sb, n)
	return sb.String()
}

const (
	verticalPipe             = '│'
	horizontalPipe           = '─'
	cornerPipe               = '└'
	verticalPipeWithOffshoot = '├'
)

func printAsTreeHelper(sb *strings.Builder, n *node) string {
	for _, child := range n.children {
		for _, parent := range child.findParents() {
			var pipe rune
			if parent.isLastChild() {
				pipe = ' '
			} else {
				pipe = verticalPipe
			}
			sb.WriteString(fmt.Sprintf("%c%s", pipe, spaces(len(parent.name)-1)))
		}

		var pipe rune
		if child.isLastChild() {
			pipe = cornerPipe
		} else {
			pipe = verticalPipeWithOffshoot
		}

		sb.WriteString(fmt.Sprintf("%c%c %s\n", pipe, horizontalPipe, child.name))
		printAsTreeHelper(sb, child)
	}
	return sb.String()
}

func spaces(n int) string {
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = ' '
	}
	return string(s)
}

func printScannerAsTree(s *bufio.Scanner) string {
	root := newNode(
		"root",
		nil,
		[]*node{},
	)
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
