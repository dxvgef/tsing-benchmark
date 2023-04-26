# tsing-benchmark
包含以下框架的路由基准测试
- github.com/dimfeld/httptreemux/v5
- github.com/dxvgef/tsing
- github.com/dxvgef/tsing/v2
- github.com/gin-gonic/gin
- github.com/julienschmidt/httprouter
- github.com/labstack/echo/v4

同时测试了启用和禁用`recover`的两种情况

#### 测试方法：
```
go test -bench=. -benchmem
```

## 测试结果
```
Benchmark_TsingV2-8                        50865             23725 ns/op               0 B/op          0 allocs/op
Benchmark_TsingV2_Recover-8                48708             24582 ns/op               0 B/op          0 allocs/op
Benchmark_TsingV1-8                        48664             24875 ns/op               0 B/op          0 allocs/op
Benchmark_TsingV1_Recover-8                45986             26267 ns/op               0 B/op          0 allocs/op
Benchmark_Gin-8                            47978             24542 ns/op               0 B/op          0 allocs/op
Benchmark_Gin_Recover-8                    43753             27390 ns/op               0 B/op          0 allocs/op
Benchmark_Httprouter-8                     46738             25555 ns/op           13792 B/op        167 allocs/op
Benchmark_Httprouter_Recover-8             44786             26703 ns/op           13792 B/op        167 allocs/op
Benchmark_Echo-8                           38401             31216 ns/op               0 B/op          0 allocs/op
Benchmark_Echo_Recover-8                   28674             41750 ns/op            9748 B/op        203 allocs/op
Benchmark_HTTPTreemux-8                    15448             77755 ns/op           65857 B/op        671 allocs/op
```
