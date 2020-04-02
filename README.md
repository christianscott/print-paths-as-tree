# print-paths-as-tree

Print a list of paths from stdin as a `tree`-style tree

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
```
