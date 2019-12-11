# tsing-benchmark
Tsing框架的基准测试代码，同时加入了以下框架的对比
- github.com/julienschmidt/httprouter 
- github.com/labstack/echo
- github.com/gin-gonic/gin

#### 使用方法：
```
go test -bench=. -benchmem
```

#### 测试结果
尽可能关闭所有框架的各种功能特性，以体现框架的最高性能
```
HttpRouter: 38504 Bytes
Tsing: 49048 Bytes
Echo: 97072 Bytes
Gin: 59128 Bytes

BenchmarkHttpRouter-4              37534             32194 ns/op           13792 B/op        167 allocs/op
BenchmarkTsing-4                   31329             37699 ns/op           13794 B/op        167 allocs/op
BenchmarkEcho-4                    32444             37044 ns/op               0 B/op          0 allocs/op
BenchmarkGin-4                     29630             40147 ns/op               0 B/op          0 allocs/op
```

开启各框架的recover功能后的测试结果
```
HttpRouter: 36696 Bytes
Tsing: 49456 Bytes
Echo: 96904 Bytes
Gin: 60872 Bytes

BenchmarkHttpRouter-4              31113             37298 ns/op           13792 B/op        167 allocs/op
BenchmarkTsing-4                   25537             47324 ns/op           13794 B/op        167 allocs/op
BenchmarkEcho-4                    19062             62752 ns/op            9745 B/op        203 allocs/op
BenchmarkGin-4                     22995             52700 ns/op               0 B/op          0 allocs/op
```

#### 总结
仅以最小功能测试时，httprouter的执行效率最好，因为它实现的功能最少。echo其次，tsing和gin水平非常接近。

但是开启了recover机制，各框架的性能都有不同程度的降低，httprouter依然表现最好，tsing也有明显的影响，但是echo和gin的影响较大，echo则是降低了近一倍的性能。而且也不再是零内存分配。

在参与测试的框架里对执行效率和功能完备之间做折中选择，tsing可能最好。