# 🧐🔗 go-cyclic
<hr/>

[ [English](README.md) | 中文 ]

⚡ Go的循环依赖检测工具 ⚡

## 🤔 什么是 go-cyclic
在 Go 应用程序的开发过程中，包与包之间的循环依赖是一个常见的问题，这类情况通常会导致编译错误，具体表现为`import cycle not allowed`的提示。当项目规模膨胀，依赖关系错综复杂时，识别并解决循环引用的问题就变得更加具有挑战性，往往会耗费开发者大量的时间和精力。

正是鉴于这一痛点，`go-cyclic`工具应运而生，它的设计初衷是为了高效且精准地帮助开发者定位项目中的循环引用问题。通过智能化的分析，`go-cyclic`能够迅速揭示循环依赖的具体位置，从而极大地简化了排查过程，保证了项目的健康与可维护性，是优化大型项目结构、提升开发效率的得力助手。

以下是会出现循环依赖的示例。
```bash
# a.go                       # b.go
package a                    package b

import "b"                   import "a"

type A struct {              type B struct {
  B *b.B                       A *a.A
}                            }
```

##快速开始
安装命令
```bash
go install github.com/elza2/go-cyclic
```
运行命令
```bash
go-cyclic run --dir .
```
go-cyclic 的参数：<br/>
`--dir` 路径参数。提示：设置的目录要为 go.mod 文件所在的目录。<br/>
`--filter` (可选) 过滤参数。提示：过滤匹配的文件，不参与循环检测。多个条件使用英文逗号隔开，支持表达式，例如 `--filter *_test.go,a_test.go`

##结果展示
1. 检测正常，无循环依赖。
```bash
Success. Not circular dependence.
```

2. 检测失败，存在循环依赖。
```bash
Failed. 1 circular dependence chains were found.

┌---→    app.go
┆          ↓
┆       routes.go
┆          ↓
└---    handler.go
```


