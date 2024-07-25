# 🧐🔗 go-cyclic
<hr/>

[ English | [中文](README_zh.md) ]

⚡ Circular dependency detection tool for Go ⚡

## 🤔 What is go-cyclic?
In the development process of Go applications, cyclic dependencies between packages are a common problem. This kind of situation usually leads to compilation errors. Specifically, execute the prompt of `import cycle not allowed`. When the project scale expands and dependencies become complex, identifying and solving circular reference problems becomes more challenging, often resulting in a lot of time and effort.

It is in view of this pain point that the `go-cyclic` tool came into being. It was originally designed to help developers locate circular reference problems in projects efficiently and accurately. Through intelligent analysis, `go-cyclic` can quickly reveal the specific location of cyclic dependencies, thus greatly simplifying the troubleshooting process and ensuring the health and maintainability of the project. It is a powerful assistant for optimizing the structure of large projects and improving development efficiency. .

The following are examples of where circular dependencies can occur.
```bash
# a.go                       # b.go
package a                    package b

import "b"                   import "a"

type A struct {              type B struct {
  B *b.B                       A *a.A
}                            }
```

## Quick Start
Install command.
```bash
go install github.com/elza2/go-cyclic@latest
```
Run command.
```bash
go-cyclic run --dir .
```
Parameters of go-cyclic:<br/>
`--dir` path parameter. Tip: The set directory must be the directory where the go.mod file is located.<br/>
`--filter` (optional) filter parameters. Tip: Filter matching files and do not participate in loop detection. Multiple conditions are separated by commas and expressions are supported, such as `--filter *_test.go,a_test.go`<br/>

## Results display
1. The detection is normal and there is no circular dependency.

```bash
Success. Not circular dependence.
```

2. Detection failed, there is a circular dependency.
```bash
Failed. 1 circular dependence chains were found.

┌---→    app.go
┆          ↓
┆       routes.go
┆          ↓
└---    handler.go
```
