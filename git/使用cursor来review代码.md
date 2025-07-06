## 使用 cursor 来 review 代码

### 前言

cursor 是什么，这里不介绍了，认为看到这篇文章的人都知道 cursor 已经 cursor 的基本用法。  

我们这里主要来聊下 cursor 中更高阶一点的功能，比如如何来进行 code review 。  

### code review 

#### review 单个文件

比如这段代码找出，两个数中的最大值。  

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println(CompareNumbers(10, 100))
}

func CompareNumbers(a, b int) int {
	log.Printf("Comparing numbers: a=%f, b=%f", a, b)

	if a > b {
		log.Printf("Result: %f > %f", a, b)
		return b
	} else if a < b {
		log.Printf("Result: %f < %f", a, b)
		return b
	} else {
		log.Printf("Result: %f == %f", a, b)
		return a
	}
}
```

其中我们很明显能看到一个，当 a>b 因该返回 a 而不是 b。这里用 code 进行 review。   

针对这段代码，使用 command + k 呼出命令框，win自行百度 。  

<img src="/img/linux/cursor-1.jpg"  alt="cursor" align="center" />

accept 接收代码的修改。  

可以看到这个有问题的代码 cursor 已经帮助我们找到并且修复了。   

<img src="/img/linux/cursor-2.jpg"  alt="cursor" align="center" />

好了这是单个文件。下面我们看看在项目开发中针对我们每次的pr提交如何进行代码 review 。  

#### 针对提交进行 code review 







