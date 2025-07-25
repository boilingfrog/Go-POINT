<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [使用 cursor 来 review 代码](#%E4%BD%BF%E7%94%A8-cursor-%E6%9D%A5-review-%E4%BB%A3%E7%A0%81)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [code review](#code-review)
    - [review 单个文件](#review-%E5%8D%95%E4%B8%AA%E6%96%87%E4%BB%B6)
    - [针对提交进行 code review](#%E9%92%88%E5%AF%B9%E6%8F%90%E4%BA%A4%E8%BF%9B%E8%A1%8C-code-review)
  - [总结](#%E6%80%BB%E7%BB%93)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 使用 cursor 来 review 代码

### 前言

cursor 是什么，这里不介绍了，认为看到这篇文章的人都知道 cursor 以及 cursor 的基本用法。

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

其中我们很明显能看到一个，当 a>b 因该返回 a 而不是 b。这里用 cursor 进行 review。   

针对这段代码，使用 command + k 呼出命令框，win自行百度 。  

<img src="/img/linux/cursor-1.jpg"  alt="cursor" align="center" />

accept 接收代码的修改。  

可以看到这个有问题的代码 cursor 已经帮助我们找到并且修复了。   

<img src="/img/linux/cursor-2.jpg"  alt="cursor" align="center" />

好了这是单个文件。下面我们看看在项目开发中针对我们每次的pr提交如何进行代码 review 。  

#### 针对提交进行 code review 

好了接着刚刚的函数，来进行一步来探讨如何针对项目级别的代码提交进行 code review。  

这里先将刚刚的代码提交，然后重新切换一个分支，在分支中修改。   

在开发分支修改成功之后，提交代码。  

使用 `git diff` 对比个分治代码的差异部分，然后将禅意部分输出到一个diff文件，然后让cursor针对这个文件进行review。  

将刚刚的代码拆分到不同的文件中，然后修改代码提交。我们还把这段代码输出写错，让 cursor 帮我们进行 review 。

然后使用 git diff 对比个分治代码的差异部分 `git diff show-diff..master > code.diff`    

在 cursor 中找到这个文件，让cursor 基于 diff 文件，来进行 code review。  

<img src="/img/linux/cursor-3.jpg"  alt="cursor" align="center" />  

可以看到 cursor 已经基于 diff文件，帮我们对提交的代码进行了 review ，找出了问题点，同时也提出了修改的意见。   

总结下，使用 cursor 来 review 代码，首先需要将代码提交到某个分支，然后切换到这个分支，然后修改代码，然后提交代码，最后使用 git diff 对比两个分支的差异，将差异输出到一个文件中，然后让 cursor 基于这个文件进行 review 。

### 总结

上面整理了使用 cursor 来 review 代码的流程，当然随着ai技术的发现可能会有更好的工具和更便捷的使用方式出现，但是还是希望上面的办法对大家的工作效率和质量的提升提供帮助。  






