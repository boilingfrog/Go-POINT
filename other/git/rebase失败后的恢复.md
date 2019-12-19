## rebase失败后的恢复

### 记一次翻车现场  

记一次翻车的现场，很早之前提的PR后面由于需求的变便去忙别的事情了，等到要做这个需求的我时候，发现已经
落后版本了，并且有很多文件的冲突，然后就用rebase去拉代码解决冲突，然后解完之后推代码，但是之后发现一
个文件在解决冲突的时候丢失了。这时候去查看git提交历史，发现rebase之后找不到这个文件了。

最后如何解决呢，要是把丢失的文件在写一遍，那就真的变成了一名咸鱼了。这时候我们就需要去弄明白rebase
的原理了。


- [rebase的使用流程](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)
- [rebase使用的意义](#rebase%e4%bd%bf%e7%94%a8%e7%9a%84%e6%84%8f%e4%b9%89)

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

使用rebase能得到一个干净的，没有merge commit的线性历史树。
![Aaron Swartz](https://github.com/zhan-liz/Go-POINT/blob/master/img/rebase_1.png?raw=true)



 

### 参考
- 【Git Rebase原理以及黄金准则详解】 https://segmentfault.com/a/1190000005937408  