## rebase失败后的恢复

### 记一次翻车现场  

记一次翻车的现场，很早之前提的PR后面由于需求的变便去忙别的事情了，等到要做这个需求的我时候，发现已经
落后版本了，并且有很多文件的冲突，然后就用rebase去拉代码解决冲突，然后解完之后推代码，但是之后发现一
个文件在解决冲突的时候丢失了。这时候去查看git提交历史，发现rebase之后找不到这个文件了。

最后如何解决呢，要是把丢失的文件在写一遍，那就真的变成了一名咸鱼了。这时候我们就需要去弄明白rebase
的原理了。


- [rebase的使用流程](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)
- [rebase使用的意义](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)
- [rebase工作的原理](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)
   - [git工作的原理](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)

### rebase的使用流程

rebase就是把主分支的最新代码拉取自己当前的开发分之上，只不过使用rebase,会形成
更加干净的git线。

那么它的使用流程：


基于forked模式的开发。
````
1、forked代码仓库
2、 git clone <个⼈仓库地址>
3、 添加远程仓库
 git remote add remote <远程仓库地址>
4、 查看远程仓库版本
git remote -v
5、 rebase
git pull remote master --rebase
6、 遇到冲突
git add .
git rebase --continue
git push -f origin XXXXX
``````

### rebase使用的意义

使用rebase的提交历史
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/rebase_2.png?raw=true)
对比merge
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/rebase_3.png?raw=true)

使用rebase会得到一个干净的，线性的提交历史，没有不必要的合并。
使用merge能够保存项目完整的历史，并且避免公共分之上的commit。


### rebase工作的原理

为了弄清楚rebase的原理，首先需要弄清楚git的工作原理。

#### git工作原理

首先我们先来了解下git的模型。
首先我们可以看到在每个项目的下面都有一个.git的隐藏目录
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/rebase_4.png?raw=true)
关于git的一切都存储在这个目录里面（全局配置除外）。这里面有一些子目录和文件，
记录到了git所有的信息。
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/rebase_5.png?raw=true)
文件里面存储的都是一些配置文件，我们主要关注下这几个子目录：

-  info：这个目录不重要，里面有一个 exclude 文件和 .gitignore 文件的作用相似，区别是这个文件不
会被纳入版本控制，所以可以做一些个人配置。
-  hooks：这个目录很容易理解， 主要用来放一些 git 钩子，在指定任务触发前后做一些自定义的配置，这
是另外一个单独的话题，本文不会具体介绍。
-  objects：用于存放所有 git 中的对象，下面单独介绍。
-  logs：用于记录各个分支的移动情况，下面单独介绍。
-  refs：用于记录所有的引用，下面单独介绍。


























 

### 参考
- 【Git Rebase原理以及黄金准则详解】 https://segmentfault.com/a/1190000005937408  
- 【图片引用】https://blog.csdn.net/chenansic/article/details/44122107