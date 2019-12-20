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
   - [git对象](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)

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
文件里面存储的都是一些配置文件:

-  info：初始化时只有这个文件，用于排除提交规则，与 .gitignore 功能类似。他们的区别在
于.gitignore 这个文件本身会提交到版本库中去，用来保存的是公共需要排除的文件；而info/exclude 
这里设置的则是你自己本地需要排除的文件，他不会影响到其他人，也不会提交到版本库中去。
-  hooks：这个目录很容易理解， 主要用来放一些 git 钩子，在指定任务触发前后做一些自定义的配置，这
是另外一个单独的话题，本文不会具体介绍。
-  objects：用于存放所有 git 中的对象，里面存储所有的数据内容，下面单独介绍。
-  logs：用于记录各个分支的移动情况，下面单独介绍。
-  refs：用于记录所有的引用，下面单独介绍。

-  HEAD：文件指示目前被检出的分支
-  index：文件保存暂存区信息

### git对象

git是面向对象的!

举个栗子

假如我们init了一个新的仓库，然后提交了两个文件，那么会有那些对象呢?

````
$ git init
$ echo 'aaaaa'>a.txt
$ echo 'bbbbb'>b.txt
$ git add *.txt
````
上面提交了两个文件到了暂存区，我们了解到对象都存储在object文件夹中，那我们去到里面看下。
````
$ tree .git/objects
.git/objects
├── cc
│   └── c3e7b48da0932cc0f7c4ce7b4fd834c7032fe1
├── db
│   └── 754dbd326f1b7c530672afbbfef8d9223033b7
├── info
└── pack

````
git cat-file [-t] [-p] 号称git里面的瑞士军刀，我们来剖析下，t可以查看object的类型，-p可
以查看object储存的具体内容。
````
$ git cat-file -t ccc3e
blob
$ git cat-file -p ccc3e
aaaaa
````
可以发现object是一个blob类型的节点，内容是aaa,这就是object存储的内容

这里我们遇到了第一种Git object，blob类型，它只储存的是一个文件的内容，不
包括文件名等其他信息。然后将这些信息经过SHA1哈希算法得到对应的哈希值
ccc3e7b48da0932cc0f7c4ce7b4fd834c7032fe1 ，作为object在git中的唯一身份证。

然后我们进行一次commit

````
 git commit -am '[+] init'
[master（根提交） d290482] [+] init
 2 files changed, 2 insertions(+)
 create mode 100644 a.txt
 create mode 100644 b.txt
````
 再次查看object
 
 ````
$ tree .git/objects
.git/objects
├── 3b
│   └── e060dffd163e01ca6f44dfc6746dda5ae19d5d
├── cc
│   └── c3e7b48da0932cc0f7c4ce7b4fd834c7032fe1
├── d2
│   └── 904827c81ba3c05072d66d1dd466e01bcdee69
├── db
│   └── 754dbd326f1b7c530672afbbfef8d9223033b7
├── info
└── pack
 ````
里面变成了四个对象了

把我们‘瑞士军刀’剖析下
````
$ git cat-file -t 3be0
tree

git cat-file -p 3be0
100644 blob ccc3e7b48da0932cc0f7c4ce7b4fd834c7032fe1	a.txt
100644 blob db754dbd326f1b7c530672afbbfef8d9223033b7	b.txt
````

我们发现这里出现了另一种object类型-tree。它将当前的提交目录打了一个快照，存储了当前提交的目录结构。

再看下两一个object里面的信息
````
$ git cat-file -p d2904
tree 3be060dffd163e01ca6f44dfc6746dda5ae19d5d
author liz <rick.lizhan@gmail.com> 1576828975 +0800
committer liz <rick.lizhan@gmail.com> 1576828975 +0800
````
里面存储的是提交者的信息

那么有一个问题，一次commit生成的快照是全局的还是提交部分的呢？
我们再来看。
修改a.txt,然后提交，看看新的提交
 
 
 
 
 
 
 
 

### 参考
- 【Git Rebase原理以及黄金准则详解】 https://segmentfault.com/a/1190000005937408  
- 【图片引用】https://blog.csdn.net/chenansic/article/details/44122107