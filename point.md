## 一些实用且易遗漏的点

### map使用的时候，如果value只是作为占位符

```go
	mapp:=make(map[string]struct{},10)
```

推荐使用`struct`。  

因为空结构体变量的内存占用大小为0，而`bool`类型内存占用大小为1，这样可以更加最大化利用我们服务器的内存空间。  


### 使用goleak监测goroutine泄露

```go
import (
	"fmt"
	"sync"
	"testing"

	"go.uber.org/goleak"
)

func TestLeakWithGoleak(t *testing.T) {
	defer goleak.VerifyNone(t)
    // 具体的函数
	Goleak()
}
```

运行

```go
$ go test -v -run ^TestLeakWithGoleak$
```