# go-cyclic

<h4> Go 循环依赖检测工具 </h4>

快速开始
===============
```bash
go install github.com/elza2/go-cyclic@latest
# path 路径要设置为 go.mod 文件所在的路径.
# filter 过滤匹配的文件, 多个条件使用逗号隔开(,)
go-cyclic gocyclic --dir .path [--filter *_test.go]
```

运行测试
===============
```bash
git clone https://github.com/elza2/go-cyclic.git
# path 路径要设置为 go.mod 文件所在的路径.
go run ./main.go gocyclic --dir .path  [--filter *_test.go]
```

运行结果
===============
```bash
# success output.
Success. Not circular dependence.

# failed output.
Failed. 1 circular dependence chains were found.

┌---→    app.go
┆          ↓
┆       routes.go
┆          ↓
└---    handler.go
```


