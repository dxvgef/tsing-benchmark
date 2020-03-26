# tsing-benchmark
Tsing框架的基准测试代码，同时加入了以下框架的对比
- github.com/julienschmidt/httprouter
- github.com/labstack/echo
- github.com/gin-gonic/gin

同时测试了启用和禁用`recover`的两种情况

#### 测试方法：
```
go test -bench=. -benchmem
```
