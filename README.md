# tsing-benchmark
`Tsing v1.5.0`框架与以下框架的路由基准测试
- github.com/julienschmidt/httprouter
- github.com/labstack/echo
- github.com/gin-gonic/gin
- github.com/go-chi/chi

同时测试了启用和禁用`recover`的两种情况

#### 测试方法：
```
go test -bench=. -benchmem
```
