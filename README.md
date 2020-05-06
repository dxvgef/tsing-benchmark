# tsing-benchmark
`Tsing v1`框架与以下框架的路由基准测试
- github.com/julienschmidt/httprouter [v1.3.0]
- github.com/labstack/echo [v3.3.10]
- github.com/gin-gonic/gin [v1.6.1]
- github.com/dxvgef/tsing-gateway [v0]

同时测试了启用和禁用`recover`的两种情况（`Tsing Gateway`没有`recover`功能）

#### 测试方法：
```
go test -bench=. -benchmem
```
