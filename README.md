# print-paths-as-tree

## Usage

`print-paths-as-tree` accepts a list of paths from stdin and then prints them as a [`tree`](http://mama.indstate.edu/users/ice/tree/)-style tree.

```
$ cat << EOF | print-paths-as-tree
> dir1/one.file
> dir1/two.file
> dir2/one.file
> EOF
.
├── dir1
│   ├── one.file
│   └── two.file
└── dir2
    └── one.file

2 directories, 3 files
```

Handy for nicely presenting affected files as a tree:

```
$ git diff --name-only | print-paths-as-tree
src
├── components
│   ├── avatar.tsx
│   └── list.tsx
└── services
    └── users.ts

3 directories, 3 files
```

## Installation

Super simple with `go get`: `go get github.com/christianscott/print-paths-as-tree`

## How it works

1. Construct a tree from the paths, each segment inside the path becoming a node (i.e. the same as the file system). For example, the paths `src/one src/two` would become:
```
  src
 /   \
one  two
```
2. Perform a pre-order depth first traversal of the tree and print row for each node. The above tree would be visited in the order `src -> one -> two`.
3. To print a row, we need to print a "connector" for each ancestor of the current node, and then print the connector + the name of the final node. The type of connector depends on whether or not the node is its parents final child.
