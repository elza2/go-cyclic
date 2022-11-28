# go-cyclic

<h4> Go 循环依赖检测工具 </h4>

快速开始
===============
```bash
go install github.com/elza2/go-cyclic@latest
# path 路径要设置为 go.mod 文件所在的全路径. 
go-cyclic gocyclic --dir .path
```

运行测试
===============
```bash
git 
go run ./cmd/main.go gocyclic --dir .path
```

```bash
# success output.
Success. Not circular dependence.

# failed output.
┌---→    daji.go
┆          ↓
┆       routes.go
┆          ↓
└---    handler.go
```


